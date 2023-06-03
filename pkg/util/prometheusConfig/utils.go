package prometheusConfig

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"minik8s/pkg/util/config"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	Global struct {
		ScrapeInterval     string `yaml:"scrape_interval"`
		EvaluationInterval string `yaml:"evaluation_interval"`
	} `yaml:"global"`
	Alerting struct {
		Alertmanagers []struct {
			StaticConfigs []struct {
				Targets []string `yaml:"targets"`
			} `yaml:"static_configs"`
		} `yaml:"alertmanagers"`
	} `yaml:"alerting"`
	RuleFiles     []string `yaml:"rule_files"`
	ScrapeConfigs []struct {
		JobName       string `yaml:"job_name"`
		StaticConfigs []struct {
			Targets []string `yaml:"targets"`
		} `yaml:"static_configs"`
	} `yaml:"scrape_configs"`
}

// Reload Prometheus
func Reload() {
	reloadUrl := config.PROMETHEUS_URL + "/-/reload"

	// Set up the HTTP POST request
	body := url.Values{}
	req, err := http.NewRequest("POST", reloadUrl, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Body = ioutil.NopCloser(strings.NewReader(body.Encode()))

	// Send the HTTP request
	resp, err := http.Post(reloadUrl, "application/x-www-form-urlencoded", req.Body)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func AddTarget(ip string) {
	// Read in the YAML configuration file
	data, err := ioutil.ReadFile("prometheus.yml")
	if err != nil {
		panic(err)
	}

	// Unmarshal the YAML into a Config struct
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}

	// Add the IP address to the targets list
	cfg.ScrapeConfigs[0].StaticConfigs[0].Targets = append(cfg.ScrapeConfigs[0].StaticConfigs[0].Targets, ip)

	// Marshal the updated Config struct back into YAML
	newData, err := yaml.Marshal(&cfg)
	if err != nil {
		panic(err)
	}

	// Write the updated YAML back to the file
	err = ioutil.WriteFile("prometheus.yml", newData, 0644)
	if err != nil {
		panic(err)
	}

	Reload()
}

func DelTarget(ip string) {
	// Read in the YAML configuration file
	data, err := ioutil.ReadFile("prometheus.yml")
	if err != nil {
		panic(err)
	}

	// Unmarshal the YAML into a Config struct
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}

	// Remove the IP address from the targets list
	for i, target := range cfg.ScrapeConfigs[0].StaticConfigs[0].Targets {
		if target == ip {
			cfg.ScrapeConfigs[0].StaticConfigs[0].Targets = append(cfg.ScrapeConfigs[0].StaticConfigs[0].Targets[:i], cfg.ScrapeConfigs[0].StaticConfigs[0].Targets[i+1:]...)
			break
		}
	}

	// Marshal the updated Config struct back into YAML
	newData, err := yaml.Marshal(&cfg)
	if err != nil {
		panic(err)
	}

	// Write the updated YAML back to the file
	err = ioutil.WriteFile("prometheus.yml", newData, 0644)
	if err != nil {
		panic(err)
	}

	Reload()
}
