// main.go (Serviço B)
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func main() {
	// Configure Zipkin exporter
	exporter, err := zipkin.New("http://localhost:9411/api/v2/spans")
	if err != nil {
		log.Fatalf("Falha ao criar Zipkin : %v", err)
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)

	http.HandleFunc("/weather", weatherHandler)
	log.Println("Serviço B rodando na porta 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	cep := r.URL.Query().Get("cep")
	if cep == "" || len(cep) != 8 {
		http.Error(w, "Cep Inválido", http.StatusUnprocessableEntity)
		return
	}

	ctx := r.Context()
	tr := otel.Tracer("Servico-B")
	ctx, span := tr.Start(ctx, "weatherHandler")
	defer span.End()

	city, err := getCityByCEP(ctx, cep)
	if err != nil {
		http.Error(w, "Não foi possível achar o CEP", http.StatusNotFound)
		return
	}

	tempC, err := getTemperatureByCity(ctx, city)
	if err != nil {
		http.Error(w, "Falha ao buscar a Temperatura !", http.StatusInternalServerError)
		return
	}

	resp := WeatherResponse{
		City:  city,
		TempC: tempC,
		TempF: tempC*1.8 + 32,
		TempK: tempC + 273,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func getCityByCEP(ctx context.Context, cep string) (string, error) {
	tr := otel.Tracer("Servico-B")
	ctx, span := tr.Start(ctx, "getCityByCEP")
	defer span.End()

	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if city, ok := result["localidade"].(string); ok {
		return city, nil
	}

	return "", fmt.Errorf("Cidade não encontrada para o CEP %s", cep)
}

func getTemperatureByCity(ctx context.Context, city string) (float64, error) {
	tr := otel.Tracer("Servico-B")
	ctx, span := tr.Start(ctx, "getTemperatureByCity")
	defer span.End()

	encodedCity := url.QueryEscape(city)

	apiKey := "569bd7c690564646842141230250601"
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, encodedCity)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("Falha para criar a requisição: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Falha ao montar os dados da temperatura: %w", err)
	}
	//DeferLog to serverB
	defer func() {
		log.Println("Span finalizado em Serviço B: getTemperatureByCity")
		span.End()
	}()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if current, ok := result["current"].(map[string]interface{}); ok {
		if tempC, ok := current["temp_c"].(float64); ok {
			return tempC, nil
		}
	}

	return 0, fmt.Errorf("temperature not found for city %s", city)
}
