package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	// Inicializar contexto e cliente
	godotenv.Load(".env")
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("API Key não encontrada. Verifique o arquivo .env.")
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatalf("Erro ao criar o cliente: %v", err)
	}
	defer client.Close()

	// Fazer upload da imagem
	imgURI, err := uploadImage(ctx, client, "MicrosoftTeams-image 2.png")
	if err != nil {
		log.Fatalf("Erro ao fazer upload da imagem: %v", err)
	}

	// Gerar resumo da imagem
	err = generateImageSummary(ctx, client, imgURI)
	if err != nil {
		log.Fatalf("Erro ao gerar resumo da imagem: %v", err)
	}
}

// Função para fazer upload da imagem
func uploadImage(ctx context.Context, client *genai.Client, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	opts := genai.UploadFileOptions{DisplayName: "Jetpack drawing"}
	img, err := client.UploadFile(ctx, "", file, &opts)
	if err != nil {
		return "", fmt.Errorf("erro ao fazer upload do arquivo: %v", err)
	}

	fmt.Printf("Arquivo %s enviado com sucesso: %q\n", img.DisplayName, img.URI)
	return img.URI, nil
}

// Função para gerar o resumo da imagem
func generateImageSummary(ctx context.Context, client *genai.Client, imgURI string) error {
	// Escolher o modelo generativo (exemplo: "gemini-1.5-pro")
	model := client.GenerativeModel("gemini-1.5-pro")

	// Criar o prompt com a URI da imagem e o texto
	prompt := []genai.Part{
		genai.FileData{URI: imgURI},
		genai.Text("Faça um resumo sobre esta imagem"),
	}

	// Gerar o conteúdo usando o modelo
	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return fmt.Errorf("erro ao gerar conteúdo: %v", err)
	}

	// Exibir o resumo gerado
	for _, c := range resp.Candidates {
		if c.Content != nil {
			fmt.Println("Resumo da imagem:", *c.Content)
		}
	}
	return nil
}
