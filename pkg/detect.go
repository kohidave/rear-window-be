package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rekognition"
)

type DetectService struct {
	rkgClient *rekognition.Rekognition
}

type DetectionResult struct {
	Scenes  []*DetectionScene
	Objects []*DetectionObject
}

type DetectionScene struct {
	Name              string
	Probability       float64
	OtherDescriptions []string
}

type DetectionObject struct {
	Name              string
	Probability       float64
	OtherDescriptions []string
	Locations         []*BoundingBox
}

type BoundingBox struct {
	// Bounding Box
	Width  float64
	Height float64
	Left   float64
	Top    float64
}

func NewDetectService() *DetectService {
	sess := session.New(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	svc := rekognition.New(sess)

	return &DetectService{
		rkgClient: svc,
	}
}

func (d *DetectService) DetectObjects(imageBytes *[]byte) (*DetectionResult, error) {
	result, err := d.rkgClient.DetectLabels(&rekognition.DetectLabelsInput{
		Image: &rekognition.Image{
			Bytes: *imageBytes,
		},
	})

	if err != nil {
		return nil, err
	}

	detections := DetectionResult{}
	for _, label := range result.Labels {
		if len(label.Instances) == 0 {
			// Scene
			scene := DetectionScene{}
			scene.Name = *label.Name
			scene.Probability = *label.Confidence
			for _, parent := range label.Parents {
				scene.OtherDescriptions = append(scene.OtherDescriptions, *parent.Name)
			}
			detections.Scenes = append(detections.Scenes, &scene)
		} else {
			obj := DetectionObject{}
			obj.Name = *label.Name
			obj.Probability = *label.Confidence
			for _, parent := range label.Parents {
				obj.OtherDescriptions = append(obj.OtherDescriptions, *parent.Name)
			}
			for _, instance := range label.Instances {
				obj.Locations = append(obj.Locations, &BoundingBox{
					Top:    *instance.BoundingBox.Top,
					Left:   *instance.BoundingBox.Left,
					Height: *instance.BoundingBox.Height,
					Width:  *instance.BoundingBox.Width,
				})
			}
			detections.Objects = append(detections.Objects, &obj)

		}
	}
	return &detections, nil
}
