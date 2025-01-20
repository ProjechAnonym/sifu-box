package singbox

import (
	"sifu-box/models"

	"go.uber.org/zap"
)

func addURLTestOutbound(outbounds []models.Outbound, tags []string, logger *zap.Logger) ([]models.Outbound, []string, error){
	var urlTest models.URLTest
	URLTestMap := map[string]interface{}{"type":"urltest", "interrupt_exist_connections":false, "tag":"auto", "outbounds": tags}
	var outbound models.Outbound = &urlTest
	outbound, err := outbound.Transform(URLTestMap, logger)
	if err != nil {
		return nil, nil, err
	} 
	tags = append(tags, outbound.GetTag())
	outbounds = append(outbounds, outbound)
	return outbounds, tags, nil
}