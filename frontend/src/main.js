import "./style.css";
import "./app.css";

import { GetTasks } from "../wailsjs/go/main/App";

const list = document.getElementById("task-list");
const status = document.getElementById("status");

async function loadTasks() {
  status.textContent = "Loading...";

  try {
    const tasks = await GetTasks();

    list.replaceChildren();

    for (const task of tasks) {
      const item = document.createElement("li");
      item.textContent = task.title;
      list.appendChild(item);
    }

    status.textContent = `${tasks.length} tasks.`;
  } catch (error) {
    console.error(error);
    status.style.color = "red";
    status.textConent = "Failed to load tasks!";
  }
}
loadTasks();
