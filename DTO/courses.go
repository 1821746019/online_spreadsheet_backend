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
	ClassName string `json:"className"`
}

type ViewCourseResponse struct {
	Cells []CourseCell `json:"cells"`
}
