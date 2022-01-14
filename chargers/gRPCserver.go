package chargers

import (
	"context"

	"github.com/RSOam/chargers-service/proto"
	gt "github.com/go-kit/kit/transport/grpc"
	"github.com/go-kit/log"
)

type gRPCServer struct {
	createCharger gt.Handler
	updateCharger gt.Handler
	deleteCharger gt.Handler
	getCharger    gt.Handler
	getChargers   gt.Handler
}

func NewGRPCServer(endpoints Endpoints, logger log.Logger) proto.ChargersServiceServer {
	return &gRPCServer{
		createCharger: gt.NewServer(
			endpoints.CreateCharger,
			decodeCreateChargerRequestG,
			encodeCreateChargerResponse,
		),
		deleteCharger: gt.NewServer(
			endpoints.DeleteCharger,
			decodeDeleteChargerRequestG,
			encodeDeleteChargerResponse,
		),
		getCharger: gt.NewServer(
			endpoints.GetCharger,
			decodeGetChargerRequestG,
			encodeGetChargerResponse,
		),
		getChargers: gt.NewServer(
			endpoints.GetChargers,
			decodeGetChargersRequestG,
			encodeGetChargersResponse,
		),
		updateCharger: gt.NewServer(
			endpoints.UpdateCharger,
			decodeUpdateChargerRequestG,
			encodeUpdateChargerResponse,
		),
	}
}

func (s *gRPCServer) CreateCharger(ctx context.Context, req *proto.CreateChargerRequest) (*proto.CreateChargerResponse, error) {
	_, resp, err := s.createCharger.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.CreateChargerResponse), nil
}
func (s *gRPCServer) DeleteCharger(ctx context.Context, req *proto.DeleteChargerRequest) (*proto.DeleteChargerResponse, error) {
	_, resp, err := s.deleteCharger.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.DeleteChargerResponse), nil
}
func (s *gRPCServer) GetChargers(ctx context.Context, req *proto.GetChargersRequest) (*proto.GetChargersResponse, error) {
	_, resp, err := s.getChargers.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.GetChargersResponse), nil
}
func (s *gRPCServer) GetCharger(ctx context.Context, req *proto.GetChargerRequest) (*proto.GetChargerResponse, error) {
	_, resp, err := s.getCharger.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.GetChargerResponse), nil
}
func (s *gRPCServer) UpdateCharger(ctx context.Context, req *proto.UpdateChargerRequest) (*proto.UpdateChargerResponse, error) {
	_, resp, err := s.updateCharger.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*proto.UpdateChargerResponse), nil
}
func encodeCreateChargerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(CreateChargerResponse)
	return &proto.CreateChargerResponse{Status: resp.Status}, nil
}
func encodeDeleteChargerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(DeleteChargerResponse)
	return &proto.DeleteChargerResponse{Status: resp.Status}, nil
}
func encodeUpdateChargerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(UpdateChargerResponse)
	return &proto.UpdateChargerResponse{Status: resp.Status}, nil
}
func encodeGetChargerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(GetChargerResponse)
	loc := &proto.Location{}
	loc.Latitude = float32(resp.Location.Latitude)
	loc.Longitude = float32(resp.Location.Longitude)
	rat := []*proto.Rating{}
	for _, e := range resp.Ratings {
		rat = append(rat, &proto.Rating{
			Id:        e.ID.Hex(),
			ChargerID: e.ChargerID.Hex(),
			UserID:    e.UserID.Hex(),
			Rating:    int64(e.Rating),
			Created:   e.Created,
			Modified:  e.Modified,
		})
	}
	com := []*proto.Comment{}
	for _, e := range resp.Comments {
		com = append(com, &proto.Comment{
			Id:        e.ID.Hex(),
			ChargerID: e.ChargerID.Hex(),
			UserID:    e.UserID.Hex(),
			Text:      e.Text,
			Created:   e.Created,
			Modified:  e.Modified,
		})
	}
	res := []*proto.Reservation{}
	for _, e := range resp.Reservations {
		res = append(res, &proto.Reservation{
			Id:        e.ID.Hex(),
			ChargerID: e.ChargerID.Hex(),
			UserID:    e.UserID.Hex(),
			From:      e.From,
			To:        e.From,
			Created:   e.Created,
			Modified:  e.Modified,
		})
	}
	wea := *&proto.Weather{}
	wea.Weather = resp.Weather.Weather
	wea.Icon = resp.Weather.Icon
	return &proto.GetChargerResponse{Name: resp.Name, Location: loc, Ratings: rat, Comments: com, Reservations: res, ClosestCity: resp.ClosestCity}, nil
}
func encodeGetChargersResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(GetChargersResponse)
	chars := []*proto.Charger{}
	for _, e := range resp.Chargers {
		chars = append(chars, charToPchar(e))
	}
	return &proto.GetChargersResponse{Chargers: chars}, nil
}
func charToPchar(resp Charger) *proto.Charger {
	char := &proto.Charger{}
	loc := &proto.Location{}
	loc.Latitude = float32(resp.Location.Latitude)
	loc.Longitude = float32(resp.Location.Longitude)
	char.AverageRating = float32(resp.AverageRating)
	rat := []*proto.Rating{}
	com := []*proto.Comment{}
	res := []*proto.Reservation{}
	char.Name = resp.Name
	char.Location = loc
	char.Ratings = rat
	char.Comments = com
	char.Reservations = res
	return char
}

func decodeCreateChargerRequestG(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.CreateChargerRequest)
	loc := Location{}
	loc.Latitude = float64(req.Location.Latitude)
	loc.Longitude = float64(req.Location.Longitude)
	return CreateChargerRequest{Name: req.Name, Location: loc}, nil
}
func decodeUpdateChargerRequestG(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.UpdateChargerRequest)
	loc := Location{}
	loc.Latitude = float64(req.Location.Latitude)
	loc.Longitude = float64(req.Location.Longitude)
	return UpdateChargerRequest{Id: req.Id, Name: req.Name, Location: loc, AverageRating: float64(req.AverageRating)}, nil
}
func decodeDeleteChargerRequestG(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.DeleteChargerRequest)
	return DeleteChargerRequest{Id: req.Id}, nil
}
func decodeGetChargerRequestG(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*proto.GetChargerRequest)
	return GetChargerRequest{Id: req.Id}, nil
}
func decodeGetChargersRequestG(_ context.Context, request interface{}) (interface{}, error) {
	_ = request.(*proto.GetChargersRequest)
	return GetChargersRequest{}, nil
}
