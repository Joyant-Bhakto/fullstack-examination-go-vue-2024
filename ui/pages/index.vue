<template>
  <div class="todo-main">
    <h1>TODOリスト</h1>

    <div v-if="statusMessage" class="status-message">{{ statusMessage }}</div>

    <div class="input-group">
      <input v-model="newTask" placeholder="新しいタスクを入力" @keyup.enter="addTodo" />
      <button @click="addTodo">追加</button>
    </div>

    <div class="search-group">
      <input v-model="searchQuery" placeholder="タスクを検索" @input="fetchTodos" />
      <select v-model="searchStatus" @change="fetchTodos">
        <option value="">すべて</option>
        <option value="created">未完了</option>
        <option value="done">完了</option>
      </select>
    </div>

    <h2>未完了のタスク</h2>
    <div v-if="todos_incomplete.length > 0">
      <div v-for="todo in todos_incomplete" :key="todo.ID" class="todo-item">
        <input
            v-if="todo.isEditing"
            v-model="todo.Task"
            class="edit-input"
            @blur="editTodo(todo)"
            @keyup.enter="editTodo(todo)"
        />
        <span v-else :class="{ 'done-task': todo.Status === 'done' }" @click="enableEdit(todo)">
          {{ todo.Task }}
        </span>
        <div class="buttons">
          <button :class="{ 'done': todo.Status === 'done' }" @click="updateStatus(todo)">
            ✔️
          </button>
          <button class="delete-button" @click="deleteTodo(todo.ID)">🗑️</button>
        </div>
      </div>
    </div>
    <div v-else>
      <p>未完了のタスクがありません。</p>
    </div>

    <h2>完了したタスク</h2>
    <div v-if="todos_completed.length > 0">
      <div v-for="todo in todos_completed" :key="todo.ID" class="todo-item">
        <input
            v-if="todo.isEditing"
            v-model="todo.Task"
            class="edit-input"
            @blur="editTodo(todo)"
            @keyup.enter="editTodo(todo)"
        />
        <span v-else :class="{ 'done-task': todo.Status === 'done' }" @click="enableEdit(todo)">
          {{ todo.Task }}
        </span>
        <div class="buttons">
          <button :class="{ 'done': todo.Status === 'done' }" @click="updateStatus(todo)">
            ✔️
          </button>
          <button class="delete-button" @click="deleteTodo(todo.ID)">🗑️</button>
        </div>
      </div>
    </div>
    <div v-else>
      <p>完了したタスクがありません。</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      newTask: '',
      todos_completed: [],
      todos_incomplete: [],
      statusMessage: '',
      searchQuery: '', // search input for task
      searchStatus: '', // dropdown for status
    };
  },
  mounted() {
    this.fetchTodos();
  },
  methods: {
    async fetchTodos() {
      try {
        let url = `/api/v1/todos?task=${this.searchQuery}&status=${this.searchStatus}`;
        const response = await fetch(url);
        if (!response.ok) throw new Error(`Failed to get todo list. statusCode: ${response.status}`);
        const data = await response.json();
        this.todos_completed = data.data.completedTasks;
        this.todos_incomplete = data.data.incompleteTasks;
      } catch (error) {
        console.error(error);
        this.statusMessage = 'タスクの取得に失敗しました';
      }
    },
    async addTodo() {
      if (this.newTask.trim() === '') return;
      try {
        const response = await fetch('/api/v1/todos', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            task: this.newTask,
            Status: 'created'
          })
        });

        if (!response.ok) throw new Error(`Failed to create todo. statusCode: ${response.status}`);
        const data = await response.json();
        this.fetchTodos(); // Refresh list after adding
        this.newTask = '';
        this.statusMessage = 'タスクが追加されました';
      } catch (error) {
        console.error('Error creating todo:', error);
        this.statusMessage = 'タスクの作成に失敗しました';
      }
    },
    enableEdit(todo) {
      todo.isEditing = true;
    },
    async editTodo(todo) {
      todo.isEditing = false;
      try {
        const response = await fetch(`/api/v1/todos/${todo.ID}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            task: todo.Task
          })
        });

        if (!response.ok) throw new Error(`Failed to edit todo. statusCode: ${response.status}`);
        this.statusMessage = 'タスクが編集されました';
        this.fetchTodos(); // Refresh list after editing
      } catch (error) {
        console.error('Error editing todo:', error);
        this.statusMessage = 'タスクの編集に失敗しました';
      }
    },
    async updateStatus(todo) {
      try {
        const response = await fetch(`/api/v1/todos/${todo.ID}`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            Status: todo.Status === 'done' ? 'created' : 'done'
          })
        });

        if (!response.ok) throw new Error(`Failed to update todo Status. statusCode: ${response.status}`);
        this.statusMessage = 'タスクのステータスが変更されました';
        this.fetchTodos(); // Refresh list after status change
      } catch (error) {
        console.error('Error updating todo Status:', error);
        this.statusMessage = 'ステータスの更新に失敗しました';
      }
    },
    async deleteTodo(id) {
      try {
        const response = await fetch(`/api/v1/todos/${id}`, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
          },
        });

        if (!response.ok) throw new Error(`Failed to delete todo. statusCode: ${response.status}`);
        this.statusMessage = 'タスクが削除されました';
        this.fetchTodos(); // Refresh list after deletion
      } catch (error) {
        console.error('Error deleting todo:', error);
        this.statusMessage = 'タスクの削除に失敗しました';
      }
    }
  }
};
</script>

<style scoped>
.todo-main {
  max-width: 400px;
  margin: 20px auto;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  background-color: #fff;
}

.input-group {
  display: flex;
  margin-bottom: 20px;
}

input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  margin-right: 10px;
}

button {
  padding: 10px;
  border: none;
  background-color: #333;
  color: #fff;
  border-radius: 4px;
  cursor: pointer;
}

.todo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.buttons button {
  background-color: #f0f0f0;
  color: #333;
  margin-left: 5px;
  border-radius: 4px;
  padding: 5px 10px;
  transition: background-color 0.3s, color 0.3s;
}

.buttons button.done {
  background-color: #333;
  color: #fff;
}

.buttons button.done::before {
  color: #fff;
}

.buttons button.delete-button {
  color: white;
}

.status-message {
  margin-bottom: 20px;
  padding: 10px;
  background-color: #f0f0f0;
  border-radius: 4px;
}

.done-task {
  text-decoration: line-through;
}

.edit-input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
.search-group {
  display: flex;
  margin-bottom: 20px;
  gap: 10px;
}

.search-group input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.search-group select {
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: #fff;
}

</style>
