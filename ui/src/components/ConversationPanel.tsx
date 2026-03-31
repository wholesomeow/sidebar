import { FolderPlus, MessageSquarePlus, FilePlusCorner, ChevronDown, ChevronRight, Settings, CircleQuestionMark } from "lucide-react";
import { useState, useEffect } from "react";

interface ConversationListItem {
  name: string;
  id: string;
  type: "file" | "folder";
  path: string;
}

interface ConversationPanelProps {
  onSelectConversation: (id: string) => void;
}

interface FolderContentsProps {
  path: string;
  expandedFolders: Set<string>;
  toggleFolder: (path: string) => void;
  onSelectConversation: (path: string) => void;
}

function FolderContents({ path, expandedFolders, toggleFolder, onSelectConversation }: FolderContentsProps) {
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
        <li
          key={item.path}
          className="py-1 pl-2 text-sm rounded hover:bg-gray-100"
          onClick={() => item.type === "file" && onSelectConversation(item.id)}
        >
          {item.type === "folder" ? (
            <div className="flex items-center gap-2" onClick={() => toggleFolder(item.path)}>
              {expandedFolders.has(item.path) ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
              <span>{item.name}</span>
            </div>
          ) : (
            <div className="flex items-center gap-2 ml-2">
              <span>{item.name}</span>
            </div>
          )}

          {item.type === "folder" && expandedFolders.has(item.path) && (
            <FolderContents
              path={item.path}
              expandedFolders={expandedFolders}
              toggleFolder={toggleFolder}
              onSelectConversation={onSelectConversation}
            />
          )}
        </li>
      ))}
    </ul>
  );
}

export default function ConversationPanel({ onSelectConversation }: ConversationPanelProps) {
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
    <div className="w-64 bg-gray-100 border-r border-gray-300 flex flex-col h-full">

      {/* Header */}
      <div className="p-2">
        <div className="flex items-center justify-center mb-2">

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
      </div>

      {/* Scrollable Conversation List */}
      <div className="flex-1 overflow-y-auto px-2">
        <ul>
          {items.map(item => (
            <li
              key={item.path}
              className="py-1 pl-2 text-sm rounded hover:bg-gray-200"
            >
              {item.type === "folder" ? (
                <div className="flex items-center gap-2" onClick={() => toggleFolder(item.path)}>
                  {expandedFolders.has(item.path) ? (
                    <ChevronDown size={16} />
                  ) : (
                    <ChevronRight size={16} />
                  )}
                  <span>{item.name}</span>
                </div>
              ) : (
                <div className="flex items-center gap-2">
                  <span className="ml-4"></span>
                  <span>{item.name}</span>
                </div>
              )}

              {item.type === "folder" && expandedFolders.has(item.path) && (
              <FolderContents
                path={item.path}
                expandedFolders={expandedFolders}
                toggleFolder={toggleFolder}
                onSelectConversation={onSelectConversation}
              />
            )}
            </li>
          ))}
        </ul>
      </div> 
      
      {/* Footer Buttons */}
      <div className="flex items-center justify-end mb-2">
        <button className="px-1" title="About">
          <CircleQuestionMark size={18} />
        </button>
        <button className="px-4" title="Settings">
          <Settings size={18} />
        </button>
      </div>
    </div>
  );
}
