package service

import (
	"encoding/json"
	"log"
	"netease-kit/nemo/internal/dto"
	"netease-kit/nemo/internal/model"
	"netease-kit/nemo/internal/repo"
	"strconv"
)

type GameReportService struct {
	GameReportRepo *repo.GameReportRepository
	GameRecordRepo *repo.GameRecordRepository
	GameMemberRepo *repo.GameMemberRepository
}

func NewGameReportService() *GameReportService {
	return &GameReportService{
		GameReportRepo: repo.NewGameReportRepository(),
		GameRecordRepo: repo.NewGameRecordRepository(),
		GameMemberRepo: repo.NewGameMemberRepository(),
	}
}

func (s *GameReportService) ReportGameInfo(reportGameInfo string) {
	if reportGameInfo == "" {
		log.Println("GameReportService: Report content is empty")
		return
	}

	var jsonObject map[string]interface{}
	if err := json.Unmarshal([]byte(reportGameInfo), &jsonObject); err != nil {
		log.Printf("GameReportService: JSON parse error: %v", err)
		return
	}

	reportType, ok := jsonObject["report_type"].(string)
	if !ok {
		log.Println("GameReportService: Missing report_type")
		return
	}

	reportMsg, ok := jsonObject["report_msg"].(map[string]interface{})
	if !ok {
		log.Println("GameReportService: Missing report_msg")
		return
	}

	reportMsgBytes, _ := json.Marshal(reportMsg)

	switch reportType {
	case "game_start":
		log.Printf("GameReportService: Game Start Report: %s", reportGameInfo)
		var gameStartDto dto.GameStartDto
		if err := json.Unmarshal(reportMsgBytes, &gameStartDto); err != nil {
			log.Printf("GameReportService: GameStartDto parse error: %v", err)
			return
		}

		if gameStartDto.ReportGameInfoKey != "" {
			gameRecordId, _ := strconv.ParseUint(gameStartDto.ReportGameInfoKey, 10, 64)
			report := &model.GameReport{
				GameRecordId: gameRecordId,
				ReportMsg:    reportGameInfo,
			}
			s.GameReportRepo.Insert(report)
		}

	case "game_settle":
		log.Printf("GameReportService: Game Settle Report: %s", reportGameInfo)
		var gameSettleDto dto.GameSettleDto
		if err := json.Unmarshal(reportMsgBytes, &gameSettleDto); err != nil {
			log.Printf("GameReportService: GameSettleDto parse error: %v", err)
			return
		}

		if gameSettleDto.ReportGameInfoKey != "" {
			gameRecordId, _ := strconv.ParseUint(gameSettleDto.ReportGameInfoKey, 10, 64)
			report := &model.GameReport{
				GameRecordId: gameRecordId,
				ReportMsg:    reportGameInfo,
			}
			s.GameReportRepo.Insert(report)

			// Clean up game record and members
			s.GameMemberRepo.DeleteByGameRecordId(gameRecordId)
			s.GameRecordRepo.Delete(gameRecordId)
		}

	default:
		log.Printf("GameReportService: Unknown report type: %s", reportType)
	}
}
