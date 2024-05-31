package handlers

import (
	"blockchain.api/internal/services"

	"github.com/gofiber/fiber/v2"
)

// BlockAddRequest é a estrutura da solicitação para adicionar um bloco
type BlockAddRequest struct {
    Data string `json:"data"`
}

// BlockAddHandler é o manipulador para adicionar um bloco à blockchain
// BlockAddHandler é o manipulador para adicionar um bloco à blockchain
func BlockAddHandler(blockchainService *services.BlockchainService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Decodifica a solicitação JSON
        var req BlockAddRequest
        if err := c.BodyParser(&req); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Bad Request"})
        }

        // Chame o método AddBlock no serviço BlockchainService passando os dados
        if err := blockchainService.AddBlock(req.Data); err != nil {
            // Se ocorrer um erro ao adicionar o bloco, retorne um erro de servidor interno
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add block"})
        }

        // Retorna uma resposta de sucesso
        return c.JSON(fiber.Map{"message": "Block added successfully"})
    }
}


// BlockGetHandler é o manipulador para obter informações sobre um bloco da blockchain
// func BlockGetHandler(c *fiber.Ctx, blockchainService *services.BlockchainService) error {
//     // Aqui você obteria informações sobre um bloco específico da blockchain
//     // Por exemplo:
	
//     blockInfo := blockchainService.GetBlockInfo("blockID")


//     // Retorna as informações do bloco em formato JSON
//     // return c.JSON(blockInfo)

//     // Exemplo de resposta enquanto não há lógica implementada
//     return c.SendString("Block Get Handler")
// }
