package models_test

import (
	"os"
	"testing"

	"github.com/lwinmgmg/http_data_store/modules/models"
	"gorm.io/gorm"
)

var dMgr *models.HdsModel

func TestMain(m *testing.M) {
	models.SettingUp("test")
	models.CreateTable(&DummyTestModel{})
	dMgr = GetDummyTestModelManager()
	exitVal := m.Run()
	models.TearDown("test")
	models.DropTable(&DummyTestModel{})
	os.Exit(exitVal)
}

type Creator interface {
	Create() (any, error)
}

type DummyTestModel struct {
	gorm.Model
	Name string
	Age  int
}

type DummyTestModelRead struct {
	ID   uint
	Name string
	Age  int
}

func (mdl *DummyTestModel) GetID() uint {
	return mdl.ID
}

func (mdl *DummyTestModel) TableName() string {
	return "test_dummy"
}

func GetDummyTestModelManager() *models.HdsModel {
	return &models.HdsModel{
		Table: &DummyTestModel{},
	}
}

func CreateRecords(records ...Creator) {
	for _, v := range records {
		v.Create()
	}
}

func TestCreate(t *testing.T) {
	dummy := DummyTestModel{Name: "Lwin Mg Mg", Age: 12}
	readData := DummyTestModelRead{}
	if err := dMgr.Create(&dummy, &readData); err != nil {
		t.Error("Error on record create :", err)
	}
	if readData.Name != dummy.Name || readData.Age != dummy.Age {
		t.Errorf("Expecting : %v, Getting : %v", dummy, readData)
	}
}

func TestGetByID(t *testing.T) {
	inputData := DummyTestModel{Name: "Lwin Mg Mg1", Age: 20}
	readData := DummyTestModelRead{}
	if err := dMgr.Create(&inputData); err != nil {
		t.Error("Error on record create :", err)
	}
	if err := dMgr.GetByID(inputData.ID, &readData); err != nil {
		t.Error("Error on get by id :", err)
	}
	if readData.ID != inputData.ID {
		t.Errorf("Expecting : %v, Getting : %v", readData.ID, inputData.ID)
	}
	if readData.Name != inputData.Name {
		t.Errorf("Expecting : %v, Getting : %v", readData.Name, inputData.Name)
	}
	if readData.Age != inputData.Age {
		t.Errorf("Expecting : %v, Getting : %v", readData.Age, inputData.Age)
	}
}

func TestGetAll(t *testing.T) {
	dList := make([]DummyTestModelRead, 0, 50)
	if err := dMgr.GetAll(&dList); err != nil {
		t.Error("Error on record GetAll :", err)
	}
	if len(dList) == 0 {
		t.Error("No Data")
	}
	for _, v := range dList {
		if v.Name == "" {
			t.Error("Getting name as empty string")
		}
		if v.Age == 0 {
			t.Error("Getting zero age")
		}
		if v.ID == 0 {
			t.Error("Getting zero error")
		}
	}
}

func TestGetByIDs(t *testing.T) {
	dList := make([]DummyTestModelRead, 0, 50)
	if err := dMgr.GetAll(&dList); err != nil {
		t.Error("Error on record GetAll :", err)
	}
	ids := make([]uint, len(dList))
	for _, v := range dList {
		ids = append(ids, v.ID)
	}
	getByIDsData := make([]DummyTestModelRead, 0, len(dList))
	dMgr.GetByIDs(ids, &getByIDsData)
	for i := 0; i < len(dList); i++ {
		if dList[i].ID != getByIDsData[i].ID {
			t.Errorf("Expecting : %v, Getting : %v", dList[i].ID, getByIDsData[i].ID)
		}
		if dList[i].Name != getByIDsData[i].Name {
			t.Errorf("Expecting : %v, Getting : %v", dList[i].Name, getByIDsData[i].Name)
		}
		if dList[i].Age != getByIDsData[i].Age {
			t.Errorf("Expecting : %v, Getting : %v", dList[i].Age, getByIDsData[i].Age)
		}
	}
}
