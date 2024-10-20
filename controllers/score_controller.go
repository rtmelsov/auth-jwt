package controllers

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-keycloak-jwt/models"
	"log"
	"net/http"
)

const GetScoreCardsXml = `
<S:Envelope
	xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
	<S:Body>
		<GetScoreCardsResponse
			xmlns="http://score.ws.creditinfo.com/">
			<return>
				<attributes>
					<name>IIN</name>
				</attributes>
				<name> BehaviorScoring</name>
			</return>
		</GetScoreCardsResponse>
	</S:Body>
</S:Envelope>
`

const ScoreXml = `
<S:Envelope
	xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
	<S:Body>
		<ScoreResponse
			xmlns="http://score.ws.creditinfo.com/">
			<return>
				<IdQuery>3619</IdQuery>
				<ErrorCode>0</ErrorCode>
				<ErrorString>Выполнено успешно</ErrorString>
				<Score>73.0</Score>
				<OneYearProbabilityOfDefault>2% - 3%</OneYearProbabilityOfDefault>
				<RiskGrade/>
				<ScoreByML>444.0</ScoreByML>
				<OneYearProbabilityOfDefaultByML>10% - 15%</OneYearProbabilityOfDefaultByML>
				<RiskGradeByML>B1</RiskGradeByML>
				<Causes>
					<name>Test</name>
					<causeText>Пример сообщения о низком балле</causeText>
				</Causes>
			</return>
		</ScoreResponse>
	</S:Body>
</S:Envelope>

`

// @Summary Get score cards static
// @Description Получить статические карточки с результатами
// @Tags scores
// @Produce json
// @Accept json
// @Produce json
// @Param login body models.ScoreCardsRequest true "ScoreCardsRequest"
// @Success 200 {object} models.ScoreCardsRequest
// @Failure 404 {object} map[string]string "ScoreCardsRequest not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /get-score-cards [post]
func GetScoreCards(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to extract user ID from token"})
		return
	}
	// Преобразуем userID в строку
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userID type"})
		return
	}

	userName, exists := c.Get("userName")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to extract user name from token"})
		return
	}

	userNameStr, ok := userName.(string)

	var scoreCards models.ScoreCardsRequest

	// Привязка JSON-данных к структуре
	if err := c.ShouldBindJSON(&scoreCards); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем тело в JSON
	jsonBody, err := json.Marshal(scoreCards)
	if err != nil {
		log.Printf("Ошибка при преобразовании тела в JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create body JSON"})
		return
	}

	req, err := http.NewRequest("POST", "http://example.com/your-endpoint", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Ошибка при создании запроса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create body JSON"})
		return
	}
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Culture", "ru-RU")
	req.Header.Set("Password", "741102401283")
	req.Header.Set("SecurityToken", tokenString)
	req.Header.Set("UserId", userIDStr)
	req.Header.Set("UserName", userNameStr)
	req.Header.Set("Version", "1")

	fmt.Println("-----------")
	fmt.Println(tokenString)
	fmt.Println("-----------")
	fmt.Println(userIDStr)
	fmt.Println("-----------")
	fmt.Println(userNameStr)
	fmt.Println("-----------")

	// Читаем XML-ответ
	xmlResponse := []byte(GetScoreCardsXml)

	// Парсим XML-ответ
	var envelope models.ScoreCardsEnvelopeXml
	err = xml.Unmarshal(xmlResponse, &envelope)
	if err != nil {
		log.Printf("Ошибка парсинга XML: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse XML"})
		return
	}

	// Преобразуем в JSON
	jsonResponse, err := json.Marshal(envelope.Body.GetScoreCardsResponse.ReturnData)
	if err != nil {
		log.Printf("Ошибка преобразования в JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert to JSON"})
		return
	}

	// Отправляем JSON-ответ клиенту
	c.JSON(http.StatusOK, gin.H{"response": string(jsonResponse)})
}

// @Summary Get score static
// @Description Получить статические результаты
// @Tags scores
// @Produce json
// @Accept json
// @Produce json
// @Param login body models.ScoreRequest true "ScoreRequest"
// @Success 200 {object} models.ScoreRequest
// @Failure 404 {object} map[string]string "ScoreRequest not found"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Security BearerAuth
// @Router /score [post]
func PostScore(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to extract user ID from token"})
		return
	}
	// Преобразуем userID в строку
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userID type"})
		return
	}

	userName, exists := c.Get("userName")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to extract user name from token"})
		return
	}

	userNameStr, ok := userName.(string)

	var score models.ScoreRequest

	// Привязка JSON-данных к структуре
	if err := c.ShouldBindJSON(&score); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем тело в JSON
	jsonBody, err := json.Marshal(score)
	if err != nil {
		log.Printf("Ошибка при преобразовании тела в JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create body JSON"})
		return
	}

	req, err := http.NewRequest("POST", "http://example.com/your-endpoint", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Ошибка при создании запроса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create body JSON"})
		return
	}
	req.Header.Set("Content-Type", "application/xml")
	req.Header.Set("Culture", "ru-RU")
	req.Header.Set("Password", "741102401283")
	req.Header.Set("SecurityToken", tokenString)
	req.Header.Set("UserId", userIDStr)
	req.Header.Set("UserName", userNameStr)
	req.Header.Set("Version", "1")

	fmt.Println("-----------")
	fmt.Println(tokenString)
	fmt.Println("-----------")
	fmt.Println(userIDStr)
	fmt.Println("-----------")
	fmt.Println(userNameStr)
	fmt.Println("-----------")

	// Читаем XML-ответ
	xmlResponse := []byte(ScoreXml)

	// Парсим XML-ответ
	var envelope models.ScoreResponseXml
	err = xml.Unmarshal(xmlResponse, &envelope)
	if err != nil {
		log.Printf("Ошибка парсинга XML: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse XML"})
		return
	}

	// Преобразуем в JSON
	jsonResponse, err := json.Marshal(envelope)
	if err != nil {
		log.Printf("Ошибка преобразования в JSON: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert to JSON"})
		return
	}

	// Отправляем JSON-ответ клиенту
	c.JSON(http.StatusOK, gin.H{"response": string(jsonResponse)})
}
