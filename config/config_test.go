package config_test

import (
	"os"
	"testing"

	"github.com/duythinht/tg/config"
	"gopkg.in/yaml.v3"
)

func TestLoadConfigDefault(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("load config %s", err)
	}

	var defaultConfig config.Config

	f, _ := os.Open("./config.yaml")

	yaml.NewDecoder(f).Decode(&defaultConfig)

	if cfg.DB.DSN != defaultConfig.DB.DSN {
		t.Fail()
	}

	if cfg.Queue.Brokers != defaultConfig.Queue.Brokers {
		t.Logf("brokers: %s - %s\n", cfg.Queue.Brokers, defaultConfig.Queue.Brokers)
		t.Fail()
	}
	if cfg.Queue.GroupID != defaultConfig.Queue.GroupID {
		t.Logf("group: %s - %s\n", cfg.Queue.GroupID, defaultConfig.Queue.GroupID)
		t.Fail()
	}

	if cfg.Queue.Topic != defaultConfig.Queue.Topic {
		t.Logf("topic: %s - %s\n", cfg.Queue.Topic, defaultConfig.Queue.Topic)
		t.Fail()
	}
}

func TestLoadConfigENV(t *testing.T) {

	topic := "xyz"
	groupId := "any-group"
	os.Setenv("QUEUE_TOPIC", topic)
	os.Setenv("QUEUE_GROUP_ID", groupId)

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("load config %s", err)
	}

	if cfg.Queue.Topic != topic {
		t.Logf("topic is not match")
		t.Fail()
	}

	if cfg.Queue.GroupID != groupId {
		t.Logf("group_id is not match")
		t.Fail()
	}
}
