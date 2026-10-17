import "./style.css";
import "./app.css";

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

async function loadTasks() {
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

      const editButton = document.createElement("button");
      editButton.textContent = "Edit";
      editButton.addEventListener("click", async () => {
        const titleEditor = document.createElement("input");
        titleEditor.type = "text";
        titleEditor.value = task.title;

        item.replaceChild(titleEditor, taskTitle);
        titleEditor.focus();

        const saveButton = document.createElement("button");
        saveButton.textContent = "Save";

        saveButton.addEventListener("click", async () => {
          try {
            await UpdateTitle(task.id, titleEditor.value);
            await loadTasks();
          } catch (error) {
            console.error(error);
            status.textContent = `Failed to update task: ${error.message}`;
          }
        });

        item.replaceChild(saveButton, editButton);
      });

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

      item.append(checkbox, taskTitle, editButton, deleteButton);
      list.appendChild(item);
    }

    status.textContent = `${tasks.length} tasks.`;
  } catch (error) {
    console.error(error);
    status.style.color = "red";
    status.textContent = "Failed to load tasks!";
  }
  status.style.color = "";
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
