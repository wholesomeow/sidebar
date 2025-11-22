import { FolderPlus, MessageSquarePlus, FilePlusCorner, ChevronDown, ChevronRight } from "lucide-react";
import { useState, useEffect } from "react";

interface ConversationListItem {
  name: string;
  type: "file" | "folder";
  path: string;
}

interface FolderContentsProps {
  path: string;
  expandedFolders: Set<string>;
  toggleFolder: (path: string) => void;
}

function FolderContents({ path, expandedFolders, toggleFolder }: FolderContentsProps) {
  const [contents, setContents] = useState<ConversationListItem[]>([]);

  useEffect(() => {
    fetch(`http://localhost:8080/api/v1/conversations?path=${encodeURIComponent(path)}`)
      .then(res => res.json())
      .then(json => setContents(json.Data ?? []))
      .catch(err => console.error("Failed to load folder contents:", err));
  }, [path]);

  return (
    <ul className="ml-4 border-l border-gray-300 pl-2">
      {contents.map(item => (
        <li key={item.path} className="p-1 rounded hover:bg-gray-100 cursor-pointer">
          {item.type === "folder" ? (
            <div className="flex items-center gap-2" onClick={() => toggleFolder(item.path)}>
              {expandedFolders.has(item.path) ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
              <span>{item.name}</span>
            </div>
          ) : (
            <div className="flex items-center gap-2 ml-3">
              <span>{item.name}</span>
            </div>
          )}

          {item.type === "folder" && expandedFolders.has(item.path) && (
            <FolderContents
              path={item.path}
              expandedFolders={expandedFolders}
              toggleFolder={toggleFolder}
            />
          )}
        </li>
      ))}
    </ul>
  );
}

export default function ConversationPanel() {
  const [items, setItems] = useState<ConversationListItem[]>([]);
  const [expandedFolders, setExpandedFolders] = useState<Set<string>>(new Set());

  const toggleFolder = (path: string) => {
    setExpandedFolders(prev => {
      const newSet = new Set(prev);
      if (newSet.has(path)) {
        newSet.delete(path);
      } else {
        newSet.add(path);
      }
      return newSet;
    });
  };

  useEffect(() => {
    fetch("http://localhost:8080/api/v1/conversations")
      .then(async res => {
        const json = await res.json();
        console.log("Full API response:", json);
        if (!res.ok) {
          console.error("Server returned error:", json);
          return;
        }
        setItems(json.Data ?? []);
      })
      .catch(err => console.error("Failed to load conversations: ", err));
  }, []);

  return (
    <div className="w-64 bg-gray-100 border-r border-gray-300 p-2 overflow-y-auto">
      {/* Top Toolbar */}
      <div className="flex items-center justify-between mb-2">
        <h2 className="text-lg font-bold mb-4">Conversations</h2>
        <div className="flex gap-2">
          <button className="p-1 rounded hover:bg-gray-200" title="New Conversation">
            <MessageSquarePlus size={18} />
          </button>
          <button className="p-1 rounded hover:bg-gray-200" title="New Note">
            <FilePlusCorner size={18} />
          </button>
          <button className="p-1 rounded hover:bg-gray-200" title="New Folder">
            <FolderPlus size={18} />
          </button>
        </div>
      </div>

      {/* Conversation List */}
      <ul>
        {items.map(item => (
          <li key={item.path} className="p-2 rounded hover:bg-gray-200 cursor-pointer">
            {item.type === "folder" ? (
              <div className="flex items-center gap-2" onClick={() => toggleFolder(item.path)}>
                {expandedFolders.has(item.path) ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
                <span>{item.name}</span>
              </div>
            ) : (
              <div className="flex items-center gap-2 ml-6">
                <span>{item.name}</span>
              </div>
            )}

            {item.type === "folder" && expandedFolders.has(item.path) && (
              <FolderContents
                path={item.path}
                expandedFolders={expandedFolders}
                toggleFolder={toggleFolder}
              />
            )}
          </li>
        ))}
      </ul>
    </div>
  );
}
