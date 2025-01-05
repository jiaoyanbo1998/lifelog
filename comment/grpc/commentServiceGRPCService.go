package grpc

import (
	"context"
	commentv1 "lifelog-grpc/api/proto/gen/api/proto/comment/v1"
	"lifelog-grpc/comment/domain"
	"lifelog-grpc/comment/service"
)

type CommentServiceGRPCService struct {
	commentService service.CommentService
	commentv1.UnimplementedCommentServiceServer
}

func NewCommentServiceGRPCService(commentService service.CommentService) *CommentServiceGRPCService {
	return &CommentServiceGRPCService{
		commentService: commentService,
	}
}

func (c *CommentServiceGRPCService) CreateComment(ctx context.Context, request *commentv1.CreateCommentRequest) (*commentv1.CreateCommentResponse, error) {
	err := c.commentService.CreateComment(ctx, domain.CommentDomain{
		UserId:   request.GetComment().GetUserId(),
		Biz:      request.GetComment().GetBiz(),
		BizId:    request.GetComment().GetBizId(),
		Content:  request.GetComment().GetContent(),
		ParentId: request.GetComment().GetParentId(),
		RootId:   request.GetComment().GetRootId(),
	})
	if err != nil {
		return &commentv1.CreateCommentResponse{}, err
	}
	return &commentv1.CreateCommentResponse{}, nil
}

func (c *CommentServiceGRPCService) DeleteComment(ctx context.Context, request *commentv1.DeleteCommentRequest) (*commentv1.DeleteCommentResponse, error) {
	err := c.commentService.DeleteComment(ctx, request.GetId())
	if err != nil {
		return &commentv1.DeleteCommentResponse{}, err
	}
	return &commentv1.DeleteCommentResponse{}, nil
}

func (c *CommentServiceGRPCService) EditComment(ctx context.Context, request *commentv1.EditCommentRequest) (*commentv1.EditCommentResponse, error) {
	err := c.commentService.EditComment(ctx, domain.CommentDomain{
		Id:      request.GetComment().GetId(),
		Content: request.GetComment().GetContent(),
	})
	if err != nil {
		return &commentv1.EditCommentResponse{}, err
	}
	return &commentv1.EditCommentResponse{}, nil
}

func (c *CommentServiceGRPCService) FirstList(ctx context.Context, request *commentv1.FirstListRequest) (*commentv1.FirstListResponse, error) {
	list, err := c.commentService.FirstList(ctx, request.GetComment().Biz,
		request.GetComment().GetBizId(), request.GetMin())
	if err != nil {
		return &commentv1.FirstListResponse{}, err
	}
	return &commentv1.FirstListResponse{
		Comments: c.toGRPCCommentDomain(list),
	}, nil
}

func (c *CommentServiceGRPCService) EveryRootChildSonList(ctx context.Context, request *commentv1.EveryRootChildSonListRequest) (*commentv1.EveryRootChildSonListResponse, error) {
	list, err := c.commentService.EveryRootChildSonList(ctx, request.GetComment().GetId(),
		request.GetComment().GetRootId(), request.GetLimit())
	if err != nil {
		return &commentv1.EveryRootChildSonListResponse{}, err
	}
	return &commentv1.EveryRootChildSonListResponse{
		Comments: c.toGRPCCommentDomain(list),
	}, nil
}

func (c *CommentServiceGRPCService) SonList(ctx context.Context, request *commentv1.SonListRequest) (*commentv1.SonListResponse, error) {
	list, err := c.commentService.SonList(ctx, request.GetComment().GetParentId(),
		request.GetLimit(), request.GetOffset())
	if err != nil {
		return &commentv1.SonListResponse{}, err
	}
	return &commentv1.SonListResponse{
		Comments: c.toGRPCCommentDomain(list),
	}, nil
}

// 将[]domain.CommonDomain，转为[]*CommentDomain
func (c *CommentServiceGRPCService) toGRPCCommentDomain(lists []domain.CommentDomain) []*commentv1.Comment {
	cs := make([]*commentv1.Comment, 0, len(lists))
	for _, list := range lists {
		cs = append(cs, &commentv1.Comment{
			UserId:     list.UserId,
			Biz:        list.Biz,
			BizId:      list.BizId,
			Content:    list.Content,
			ParentId:   list.ParentId,
			RootId:     list.RootId,
			Id:         list.Id,
			CreateTime: list.CreateTime,
			UpdateTime: list.UpdateTime,
		})
	}
	return cs
}
