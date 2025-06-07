package DTO

type ViewCourseRequest struct {
	Week int `json:"week" binding:"required"`
}

type CourseCell struct {
	Row       int    `json:"row"`
	Col       int    `json:"col"`
	Content   string `json:"content"`
	Classroom string `json:"classroom"`
	Teacher   string `json:"teacher"`
}

type ViewCourseResponse struct {
	Cells []CourseCell `json:"cells"`
}
