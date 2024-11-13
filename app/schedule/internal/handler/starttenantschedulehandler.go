package handler

import (
	"net/http"

	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/logic"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/svc"
	"github.com/xuexihuang/go-SaasSchedule/app/schedule/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func StartTenantScheduleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StartTenantScheduleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewStartTenantScheduleLogic(r.Context(), svcCtx)
		resp, err := l.StartTenantSchedule(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
