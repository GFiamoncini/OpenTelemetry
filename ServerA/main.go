// main.go (Serviço A)
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type CEPRequest struct {
	CEP string `json:"cep"`
}

func main() {
	// Configure Zipkin exporter
	exporter, err := zipkin.New("http://zipkin:9411/api/v2/spans")
	if err != nil {
		log.Fatalf("Failed to create Zipkin exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	otel.SetTracerProvider(tp)

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	http.HandleFunc("/cep", cepHandler)
	log.Println("Serviço A rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func cepHandler(w http.ResponseWriter, r *http.Request) {
	var cepReq CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&cepReq); err != nil {
		http.Error(w, "Requisição Inválida", http.StatusBadRequest)
		return
	}

	if !isValidCEP(cepReq.CEP) {
		http.Error(w, "Cep Inválido", http.StatusUnprocessableEntity)
		return
	}

	ctx := r.Context()
	tr := otel.Tracer("Servico-A")
	ctx, span := tr.Start(ctx, "cepHandler")
	defer span.End()

	resp, err := forwardToServiceB(ctx, cepReq.CEP)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// Regex para validar cep
func isValidCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}

func forwardToServiceB(ctx context.Context, cep string) ([]byte, error) {
	tr := otel.Tracer("Servico-A")
	ctx, span := tr.Start(ctx, "forwardToServiceB")

	defer func() {
		log.Println("Span finalizado em Serviço A: forwardToServiceB")
		span.End()
	}()

	log.Printf("Span criado em Serviço A: forwardToServiceB com CEP %s", cep)

	url := fmt.Sprintf("http://server-b:8081/weather?cep=%s", cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Falha ao montar dados de temperatura: %s", resp.Status)
	}

	log.Println("Resposta recebida com sucesso do Serviço B")

	// Decodifica o corpo da resposta
	var body map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		log.Printf("Erro ao decodificar resposta do Serviço B: %v", err)
		return nil, err
	}

	log.Println("Resposta decodificada com sucesso:", body)

	// Re-encode o JSON em bytes para retornar ao chamador
	responseBytes, err := json.Marshal(body)
	if err != nil {
		log.Printf("Erro ao serializar resposta decodificada: %v", err)
		return nil, err
	}

	return responseBytes, nil
}
