package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/maveonair/gh-release-exporter/internal/config"
	"github.com/maveonair/gh-release-exporter/internal/github"
	"github.com/maveonair/gh-release-exporter/internal/metrics"
	"github.com/maveonair/gh-release-exporter/internal/releases"

	"github.com/Masterminds/semver"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)
}

func main() {
	configFilePathPtr := flag.String("config", "", "config file")
	flag.Parse()

	if *configFilePathPtr == "" {
		log.Fatal("argument -config is not set")
	}

	configuration, err := config.NewConfig(*configFilePathPtr)
	if err != nil {
		log.WithError(err).Fatal()
	}

	go metrics.Setup(configuration.ListeningAddr)

	githubClient := github.NewClient()

	for {
		checkReleasesForUpdate(githubClient, configuration.Releases)

		log.WithField("interval", configuration.Interval).Info("Sleep until next update")
		time.Sleep(configuration.Interval)
	}
}

func checkReleasesForUpdate(githubClient github.Client, releases map[string]releases.Release) {
	for key, release := range releases {
		log.WithFields(log.Fields{
			"name":               key,
			"last_known_version": release.LastKnownVersion,
		}).Info("Check latest releases")

		c, err := semver.NewConstraint(fmt.Sprintf("> %s", release.LastKnownVersion))
		if err != nil {
			metrics.IncreaseErrors()

			log.WithFields(log.Fields{
				"name":               key,
				"last_known_version": release.LastKnownVersion,
			}).Error(err)

			continue
		}

		latestReleases, err := githubClient.GetLatestReleases(release.GitHubRepo)
		if err != nil {
			metrics.IncreaseErrors()

			log.WithFields(log.Fields{
				"name":               key,
				"last_known_version": release.LastKnownVersion,
			}).Error(err)

			continue
		}

		newVersion := ""
		for _, latestRelease := range latestReleases {
			parsedTag, err := semver.NewVersion(latestRelease)
			if err != nil {
				log.WithFields(log.Fields{
					"name":               key,
					"last_known_version": release.LastKnownVersion,
					"latest_release":     latestRelease,
				}).Error(err)

				continue
			}

			if c.Check(parsedTag) {
				newVersion = latestRelease
			}
		}

		if newVersion == "" {
			metrics.SetReleaseSuccessProbe(key, 1)

			log.WithFields(log.Fields{
				"name":               key,
				"last_known_version": release.LastKnownVersion,
			}).Info("No new releases available")
		} else {
			metrics.SetReleaseSuccessProbe(key, 0)

			log.WithFields(log.Fields{
				"name":               key,
				"last_known_version": release.LastKnownVersion,
				"new_version":        newVersion,
			}).Info("New releases available")
		}
	}
}
