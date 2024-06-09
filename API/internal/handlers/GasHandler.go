package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

// const axios = require("axios");
// require("dotenv").config();

// // The chain ID of the supported network
// const chainId = 1;

// (async () => {
//   try {
//     const { data } = await axios.get(
//       `https://gas.api.infura.io/v3/${process.env.INFURA_API_KEY}/networks/${chainId}/suggestedGasFees`
//     );
//     console.log("Suggested gas fees:", data);
//   } catch (error) {
//     console.log("Server responded with:", error);
//   }
// })();

func loadEnv() string{
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	APIKEY := os.Getenv("INFURA_KEY")
	return APIKEY
}
type GasData struct {
	BaseFeeTrend              string   `json:"baseFeeTrend"`
	EstimatedBaseFee          string   `json:"estimatedBaseFee"`
	NetworkCongestion         float64  `json:"networkCongestion"`
	PriorityFeeTrend          string   `json:"priorityFeeTrend"`
	HistoricalBaseFeeRange    []string `json:"historicalBaseFeeRange"`
	HistoricalPriorityFeeRange []string `json:"historicalPriorityFeeRange"`
	LatestPriorityFeeRange    []string `json:"latestPriorityFeeRange"`
}

func BaseFeeTrendHandler(c *fiber.Ctx) error {
	data, err := fetchData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get gas data"})
	}
	return c.SendString(data.BaseFeeTrend)
}

func EstimatedBaseFeeHandler(c *fiber.Ctx) error {
	data, err := fetchData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get gas data"})
	}
	return c.SendString(data.EstimatedBaseFee)
}

func NetworkCongestionHandler(c *fiber.Ctx) error {
	data, err := fetchData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get gas data"})
	}
	return c.JSON(fiber.Map{"networkCongestion": data.NetworkCongestion})
}

func PriorityFeeTrendHandler(c *fiber.Ctx) error {
	data, err := fetchData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get gas data"})
	}
	return c.SendString(data.PriorityFeeTrend)
}

func HistoricalBaseFeeRangeHandler(c *fiber.Ctx) error {
	data, err := fetchData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get gas data"})
	}
	return c.JSON(fiber.Map{"historicalBaseFeeRange": data.HistoricalBaseFeeRange})
}

func HistoricalPriorityFeeRangeHandler(c *fiber.Ctx) error {
	data, err := fetchData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get gas data"})
	}
	return c.JSON(fiber.Map{"historicalPriorityFeeRange": data.HistoricalPriorityFeeRange})
}

func LatestPriorityFeeRangeHandler(c *fiber.Ctx) error {
	data, err := fetchData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get gas data"})
	}
	return c.JSON(fiber.Map{"latestPriorityFeeRange": data.LatestPriorityFeeRange})
}

func fetchData() (*GasData, error) {
	INFURA_API_KEY := loadEnv()
	url := "https://gas.api.infura.io/v3/" + INFURA_API_KEY + "/networks/1/suggestedGasFees"

	log.Println("URL:", url)

	getReturn, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer getReturn.Body.Close()

	if getReturn.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get gas fees: %d", getReturn.StatusCode)
	}

	var responseData GasData
	err = json.NewDecoder(getReturn.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, nil
}



func AllGasDataHandler(c *fiber.Ctx) error {
	INFURA_API_KEY := loadEnv()
	url := "https://gas.api.infura.io/v3/" + INFURA_API_KEY + "/networks/1/suggestedGasFees"

	log.Println("URL:", url)
	
	getReturn, err := http.Get(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get gas fees"})
	}
	defer getReturn.Body.Close()

	if getReturn.StatusCode != http.StatusOK {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get gas fees"})
	}

	body, err := ioutil.ReadAll(getReturn.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read response body"})
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse response body"})
	}

	return c.JSON(responseData)
}