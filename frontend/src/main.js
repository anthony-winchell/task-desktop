import "./style.css";
import "./app.css";
import "@sakun/system.css";

import {
  GetTasks,
  AddTask,
  DeleteTask,
  SetTaskCompleted,
  UpdateTitle,
} from "../wailsjs/go/main/App";

const list = document.getElementById("task-list");
const status = document.getElementById("status");
const form = document.getElementById("task-form");
const titleInput = document.getElementById("task-title");

function editTask(task, item, taskTitle, editButton, deleteButton) {
  const titleEditor = document.createElement("input");
  titleEditor.type = "text";
  titleEditor.value = task.title;

  const saveButton = document.createElement("button");
  saveButton.textContent = "Save";

  const cancelButton = document.createElement("button");
  cancelButton.textContent = "Cancel";

  function cancelEditing() {
    item.replaceChild(taskTitle, titleEditor);
    item.replaceChild(editButton, saveButton);
    item.replaceChild(deleteButton, cancelButton);
  }

  async function saveEditing() {
    saveButton.textContent = "Saving...";
    saveButton.disabled = true;

    try {
      await UpdateTitle(task.id, titleEditor.value);
      await loadTasks();
    } catch (error) {
      console.error(error);
      status.textContent = `Failed to update task: ${error}`;
    }
  }

  cancelButton.addEventListener("click", cancelEditing);
  saveButton.addEventListener("click", saveEditing);

  titleEditor.addEventListener("keydown", (event) => {
    if (event.key === "Escape") {
      cancelEditing();
    }

    if (event.key === "Enter") {
      saveEditing();
    }
  });

  item.replaceChild(titleEditor, taskTitle);
  item.replaceChild(saveButton, editButton);
  item.replaceChild(cancelButton, deleteButton);

  titleEditor.focus();
}

async function loadTasks() {
  status.style.color = "";
  status.textContent = "Loading...";

  try {
    const tasks = await GetTasks();

    list.replaceChildren();

    for (const task of tasks) {
      const item = document.createElement("li");

      const checkbox = document.createElement("input");
      checkbox.type = "checkbox";
      checkbox.checked = task.completed;

      checkbox.addEventListener("change", async () => {
        try {
          await SetTaskCompleted(task.id, checkbox.checked);
          await loadTasks();
        } catch (error) {
          console.error(error);
          status.textContent = `Failed to update task: ${error.message}`;
          await loadTasks();
        }
      });

      const taskTitle = document.createElement("span");
      taskTitle.textContent = task.title;

      const deleteButton = document.createElement("button");
      deleteButton.textContent = "Delete";
      deleteButton.addEventListener("click", async () => {
        try {
          await DeleteTask(task.id);
          await loadTasks();
        } catch (error) {
          console.error(error);
          status.textContent = `Failed to delete task: ${error.message}`;
        }
      });

      const editButton = document.createElement("button");
      editButton.textContent = "Edit";
      editButton.addEventListener("click", () => {
        editTask(task, item, taskTitle, editButton, deleteButton);
      });

      item.append(checkbox, taskTitle, editButton, deleteButton);
      list.appendChild(item);
    }

    status.textContent = `${tasks.length} tasks.`;
  } catch (error) {
    console.error(error);
    status.style.color = "red";
    status.textContent = "Failed to load tasks!";
  }
}

form.addEventListener("submit", async (event) => {
  event.preventDefault();

  const title = titleInput.value.trim();

  if (!title) {
    status.textContent = "Enter a task title.";
    return;
  }

  try {
    await AddTask(title);
    titleInput.value = "";
    await loadTasks();
  } catch (error) {
    console.error(error);
    status.textContent = `Failed to add task: ${error.message}`;
  }
});

loadTasks();
