package vault

import (
	"os"
	"testing"
)

func TestEnsureDir(t *testing.T) {
	dir := t.TempDir() + "/nasted/vault"

	if err := EnsureDir(dir); err != nil {
		t.Fatalf("Ensure() dir err: %v", err)
	}

	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("Vault dir was not created")
	}
}

func TestCreatAndList(t *testing.T) {
	dir := t.TempDir()

	f, err := Creat(dir, "todo")
	if err != nil {
		t.Fatalf("Creat() err : %v", err)
	}
	f.Close()

	note, err := List(dir)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}

	if len(note) != 1 || note[0].Name != "todo.md" {
		t.Fatalf("Wanted todo name showing %+v", note)
	}
}

func TestCreatAlredyExists(t *testing.T) {
	dir := t.TempDir()
	f, err := Creat(dir, "todo")
	if err != nil {
		t.Fatalf("Creat() err :%v", err)
	}
	f.Close()

	if _, err := Creat(dir, "todo"); err != ErrAlredyExists {
		t.Fatalf("Error: %v, want erroalredyexits", err)
	}
}

func TestOpenAndSave(t *testing.T) {
	dir := t.TempDir()

	f, err := Creat(dir, "todo")
	if err != nil {
		t.Fatalf("Create() error: %v", err)
	}
	f.Close()

	f, content, err := Open(dir, "todo.md")
	if err != nil {
		t.Fatalf("Open() error: %v", err)
	}

	if len(content) != 0 {
		t.Fatalf("open() content= %q want empty", content)
	}

	if err := Save(f, "Hello World"); err != nil {
		t.Fatalf("Save() error: %v", err)
	}

	f2, content, err := Open(dir, "todo.md")
	if err != nil {
		t.Fatalf("Open() err : %v", err)
	}
	f2.Close()

	if string(content) != "Hello World" {
		t.Fatalf("Content = %q want %q", content, "Hello World")
	}
}

func TestListEmpty(t *testing.T) {
	dir := t.TempDir()

	notes, err := List(dir)
	if err != nil{
		t.Fatalf("List() error: %v",err)
	}

	if len(notes) != 0 {
		t.Fatalf("Showing notes: %+v want empty",notes)
	}
}

func TestDelet(t *testing.T) {
	dir := t.TempDir()

	f, err := Creat(dir,"todo")
	if err != nil {
		t.Fatalf("Creat() err : %v", err)
	}	
	f.Close()

	if err := Delet(dir,"todo.md"); err != nil {
		t.Fatalf("Delet() err %v", err)
	}

	notes, err := List(dir)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(notes) != 0 {
		t.Fatalf("List = %+v expected empty after delet",notes)
	}
}

func TestRename(t *testing.T) {
	dir := t.TempDir()

	f, err := Creat(dir, "todo")
	if err != nil {
		t.Fatalf("Creat() err = %v", err)
	}
	f.Close()

	if err := Rename(dir, "todo.md", "check.md"); err != nil {
		t.Fatalf("Rename() error = %v", err)
	}

	note, err := List(dir)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}

	if len(note) != 1 || note[0].Name != "check.md" {
		t.Fatalf("LIst() %+v we wnat check.md", note)
	}
}

func TestAlredyExits(t *testing.T) {
	dir := t.TempDir()

	f, err := Creat(dir, "todo")
	if err != nil {
		t.Fatalf("Creat() error = %v", err)
	}
	f.Close()

	f, err = Creat(dir, "check")
	if err != nil {
		t.Fatalf("Creat() error = %v", err)
	}
	f.Close()

	if err = Rename(dir, "todo.md", "check.md"); err != ErrAlredyExists{
		t.Fatalf("Rename() error = %v, want alredy Exits error", err)
	}
}