package app

import (
	"testing"
)

// TestRegistryKeyMismatch verifies that "shift+tab" key is registered correctly
func TestRegistryKeyMismatch(t *testing.T) {
	factory := NewKeyCommandFactory()

	t.Run("shift+tab without spaces returns real command", func(t *testing.T) {
		// This is what Update() actually uses
		cmd := factory.CreateKeyCommand("shift+tab")

		// This should return a real command now that registry uses "shift+tab" (no spaces)
		if _, ok := cmd.(NilCommand); ok {
			t.Error("Expected real command but got NilCommand")
		}
	})

	t.Run("shift + tab with spaces returns NilCommand", func(t *testing.T) {
		// Spaces in key should not match
		cmd := factory.CreateKeyCommand("shift + tab")

		// This should return NilCommand since registry uses "shift+tab" (no spaces)
		if _, ok := cmd.(NilCommand); !ok {
			t.Error("Expected NilCommand but got a real command")
		}
	})
}

// TestContextValidationEnforced verifies that command contexts are now validated
func TestContextValidationEnforced(t *testing.T) {
	factory := NewKeyCommandFactory()

	t.Run("DKeyCommand is blocked outside kanban view", func(t *testing.T) {
		model := createTestModel()
		model.ViewManager.HideSidebar = false // Not kanban mode
		model.ViewManager.CurrentView = DetailsView

		cmd := factory.CreateKeyCommand("d")
		result := cmd.Execute(model)

		if result != nil {
			t.Error("D command should return nil when not in kanban context")
		}
	})

	t.Run("DKeyCommand allowed in kanban view", func(t *testing.T) {
		model := createTestModel()
		model.ViewManager.HideSidebar = true // Kanban mode
		model.ViewManager.CurrentView = CategoriesView

		cmd := factory.CreateKeyCommand("d")
		// Command is allowed to execute (may return nil if no valid task, but that's OK)
		_ = cmd.Execute(model)
		// Test passes if no panic occurs
	})

	t.Run("commands with empty contexts work everywhere", func(t *testing.T) {
		model := createTestModel()

		// J key has empty contexts - should work anywhere
		cmd := factory.CreateKeyCommand("j")
		cvc, ok := cmd.(*ContextValidatingCommand)
		if !ok {
			t.Fatal("Expected ContextValidatingCommand")
		}

		contexts := cvc.Contexts()
		if len(contexts) != 0 {
			t.Errorf("J key should have empty contexts, got %v", contexts)
		}

		// Should execute without blocking
		_ = cmd.Execute(model)
	})
}

// TestMessagePassingArchitecture verifies that commands use the message-passing pattern
func TestMessagePassingArchitecture(t *testing.T) {
	factory := NewKeyCommandFactory()

	t.Run("file navigation returns tea.Cmd", func(t *testing.T) {
		// Create model in details view with files
		model := createTestModel()
		model.ViewManager.CurrentView = DetailsView
		model.ViewManager.HideSidebar = false
		model.FileManager.Files = []FileInfo{
			{Name: "test1.md"},
			{Name: "test2.md"},
		}

		// J key (move down) in details view should return a command
		cmd := factory.CreateKeyCommand("j")
		result := cmd.Execute(model)

		// Navigation commands now return tea.Cmd for file navigation
		if result == nil {
			t.Error("J command should return tea.Cmd for file navigation, got nil")
		}
	})

	t.Run("kanban navigation returns tea.Cmd", func(t *testing.T) {
		model := createTestModel()
		model.ViewManager.CurrentView = CategoriesView
		model.ViewManager.HideSidebar = true // Kanban mode

		cmd := factory.CreateKeyCommand("j")
		result := cmd.Execute(model)

		if result == nil {
			t.Error("J command should return tea.Cmd for kanban navigation, got nil")
		}
	})
}

// TestDuplicateInterfaces verifies that Command interface is properly implemented
func TestDuplicateInterfaces(t *testing.T) {
	factory := NewKeyCommandFactory()
	cmd := factory.CreateKeyCommand("j")

	t.Run("Command interface via ContextValidatingCommand", func(t *testing.T) {
		// Verify it implements Command
		var _ Command = cmd
	})

	t.Run("ContextValidatingCommand wraps Command interface", func(t *testing.T) {
		cvc, ok := cmd.(*ContextValidatingCommand)
		if !ok {
			t.Fatal("Expected ContextValidatingCommand")
		}

		// Verify Description works
		description := cvc.Description()
		if description == "" {
			t.Log("Description is empty")
		}

		// Verify Contexts works
		contexts := cvc.Contexts()
		_ = contexts // Just verify it works
	})
}

// Helper function to create a minimal model for testing
func createTestModel() *Model {
	model := &Model{
		DirectoryManager: DirectoryManager{
			Companies: []Company{
				{DisplayName: "Test", FolderPathName: "test", FullPath: "/tmp/test"},
			},
			Categories:       []string{"tasks"},
			SelectedCompany:  Company{DisplayName: "Test", FolderPathName: "test", FullPath: "/tmp/test"},
			SelectedCategory: "tasks",
		},
		TaskManager: TaskManager{
			TaskCollection: TaskCollection{
				TasksByFile: make(map[string][]Task),
			},
		},
		FileManager: FileManager{
			Files: []FileInfo{
				{Name: "test.md", FullPath: "/tmp/test/tasks/test.md"},
			},
		},
		ViewManager: ViewManager{
			Ready:        true,
			CurrentView:  DetailsView,
			HideSidebar:  true, // Kanban mode
			Width:        120,
			Height:       40,
		},
	}
	return model
}
