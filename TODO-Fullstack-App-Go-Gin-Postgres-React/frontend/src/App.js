import React, { useState, useEffect } from "react";
import "./App.css";

const API_URL = "http://localhost:8081";

function App() {
  const [items, setItems] = useState([]);
  const [newItemText, setNewItemText] = useState("");
  const [editingId, setEditingId] = useState(null);
  const [editText, setEditText] = useState("");

  useEffect(() => {
    fetchItems();
  }, []);

  const fetchItems = async () => {
    try {
      const res = await fetch(`${API_URL}/items`);
      const data = await res.json();
      setItems(data.items);
    } catch (err) {
      console.error("Error fetching items:", err);
    }
  };

  const addItem = async () => {
    if (!newItemText.trim()) return;
    try {
      const res = await fetch(`${API_URL}/items`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ item: newItemText }),
      });
      const data = await res.json();
      setItems([...items, data.item]);
      setNewItemText("");
    } catch (err) {
      console.error("Error adding item:", err);
    }
  };

  const toggleDone = async (id, currentDone) => {
    try {
      const res = await fetch(`${API_URL}/items/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ done: !currentDone }),
      });
      const data = await res.json();
      setItems(items.map(item => item.id === id ? data.item : item));
    } catch (err) {
      console.error("Error updating item:", err);
    }
  };

  const deleteItem = async (id) => {
    try {
      await fetch(`${API_URL}/items/${id}`, {
        method: "DELETE",
      });
      setItems(items.filter(item => item.id !== id));
    } catch (err) {
      console.error("Error deleting item:", err);
    }
  };

  const startEdit = (item) => {
    setEditingId(item.id);
    setEditText(item.item);
  };

  const saveEdit = async (id) => {
    if (!editText.trim()) return;
    try {
      const res = await fetch(`${API_URL}/items/${id}`, {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ item: editText }),
      });
      const data = await res.json();
      setItems(items.map(item => item.id === id ? data.item : item));
      setEditingId(null);
      setEditText("");
    } catch (err) {
      console.error("Error updating item:", err);
    }
  };

  const formatDate = (dateStr) => {
    const date = new Date(dateStr);
    return date.toLocaleString();
  };

  const activeItems = items.filter(item => !item.done);
  const completedItems = items.filter(item => item.done);

  const TodoItem = ({ item }) => {
    const isEditing = editingId === item.id;
    return (
      <div className={`ListItem ${item.done ? "done" : ""}`} key={item.id}>
        <div className="ListItem-Header">
          <div
            className={`Checkbox ${item.done ? "checked" : ""}`}
            onClick={() => toggleDone(item.id, item.done)}
          >
            {item.done && "✓"}
          </div>
          
          {isEditing ? (
            <>
              <input
                className="Edit-Input"
                type="text"
                value={editText}
                onChange={(e) => setEditText(e.target.value)}
                onKeyDown={(e) => e.key === "Enter" && saveEdit(item.id)}
                autoFocus
              />
              <button className="Save-Button" onClick={() => saveEdit(item.id)}>
                Save
              </button>
            </>
          ) : (
            <>
              <div className={`Title ${item.done ? "done" : ""}`}>
                {item.item}
              </div>
              <div className="Actions">
                <button 
                  className="Action-Button Edit-Button" 
                  onClick={() => startEdit(item)}
                  title="Edit"
                >
                  ✏️
                </button>
                <button 
                  className="Action-Button Delete-Button" 
                  onClick={() => deleteItem(item.id)}
                  title="Delete"
                >
                  🗑️
                </button>
              </div>
            </>
          )}
        </div>
        <div className="Meta">
          <span>Created: {formatDate(item.created_at)}</span>
          <span>Updated: {formatDate(item.updated_at)}</span>
        </div>
      </div>
    );
  };

  return (
    <div className="App">
      <div className="Header">TODO List</div>
      
      <div className="AddBar">
        <input
          className="AddBar-Text"
          type="text"
          placeholder="Enter TODO Item"
          value={newItemText}
          onChange={(e) => setNewItemText(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && addItem()}
        />
        <button className="AddBar-Button" onClick={addItem}>
          Add Todo
        </button>
      </div>

      <div className="TodoList">
        <div>
          <div className="Section-Title">Active Todos</div>
          <div className="List">
            {activeItems.length === 0 ? (
              <div className="Empty-State">No active todos! 🎉</div>
            ) : (
              activeItems.map(item => <TodoItem key={item.id} item={item} />)
            )}
          </div>
        </div>

        {completedItems.length > 0 && (
          <div>
            <div className="Section-Title">Completed Todos</div>
            <div className="List">
              {completedItems.map(item => <TodoItem key={item.id} item={item} />)}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
