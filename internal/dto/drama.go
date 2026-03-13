package dto

// DramaCreate 创建剧集请求
type DramaCreate struct {
	DramaNo               string  `json:"drama_no" required:"true" maxLength:"32" doc:"剧集编号"`
	Title                 string  `json:"title" required:"true" maxLength:"200" doc:"剧集标题"`
	Outline               *string `json:"outline,omitempty" doc:"剧集大纲"`
	CoverImage            *string `json:"cover_image,omitempty" maxLength:"500" doc:"封面图URL"`
	Characters            string  `json:"characters" maxLength:"1000" doc:"主角id列表，逗号分隔"`
	CharacterRelationDesc string  `json:"character_relation_desc" maxLength:"1000" doc:"角色关系描述"`
	Status                *int16  `json:"status,omitempty" doc:"状态: 1-草稿 2-处理中 3-成功 4-失败"`
	TaskNo                string  `json:"task_no" maxLength:"64" doc:"任务编号"`
}

// DramaUpdate 更新剧集请求（全为可选）
type DramaUpdate struct {
	Title                 *string `json:"title,omitempty" maxLength:"200"`
	Outline               *string `json:"outline,omitempty"`
	CoverImage            *string `json:"cover_image,omitempty" maxLength:"500"`
	Characters            *string `json:"characters,omitempty" maxLength:"1000"`
	CharacterRelationDesc *string `json:"character_relation_desc,omitempty" maxLength:"1000"`
	Status                *int16  `json:"status,omitempty"`
	TaskNo                *string `json:"task_no,omitempty" maxLength:"64"`
}

// DramaResponse 剧集响应
type DramaResponse struct {
	ID                    int64   `json:"id"`
	DramaNo               string  `json:"drama_no"`
	Title                 string  `json:"title"`
	Outline               *string `json:"outline,omitempty"`
	CoverImage            *string `json:"cover_image,omitempty"`
	Characters            string  `json:"characters"`
	CharacterRelationDesc string  `json:"character_relation_desc"`
	Status                *int16  `json:"status,omitempty"`
	TaskNo                string  `json:"task_no"`
	CreateBy              string  `json:"create_by"`
	UpdateBy              string  `json:"update_by"`
	CreateAt              int64   `json:"create_at"`
	UpdateAt              int64   `json:"update_at"`
}

// DramaListResponse 剧集列表响应
type DramaListResponse struct {
	Items []DramaResponse `json:"items"`
	Total int64           `json:"total"`
}
