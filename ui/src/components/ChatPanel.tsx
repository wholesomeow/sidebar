import { useEffect, useState } from "react";

interface ChatPanelProps {
  conversationID: string | null;
}

interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
}

export default function ChatPanel({ conversationID }: ChatPanelProps) {
  const [messages, setMessages] = useState<Message[]>([]);

  useEffect(() => {
    if (!conversationID) {
      setMessages([]); // clear when no conversation selected
      return;
    }

    fetch(`http://localhost:8080/api/v1/conversations/${encodeURIComponent(conversationID)}/messages`)
      .then(async res => {
        const json = await res.json();
        console.log("Full API response:", json);
        if (!res.ok) {
          console.error("Server returned error:", json);
          return;
        }
        setMessages(Object.values(json.Data ?? []));
      })
      .catch(err => {
        console.error("Failed to load messages:", err);
        setMessages([]);
      });
  }, [conversationID]);

  return (
    <div className="flex-1 flex flex-col bg-white">
      <div className="flex-1 p-4 overflow-y-auto">
        {messages.map(msg => (
          <div key={msg.id} className={`mb-4 ${msg.role === 'user' ? 'text-right' : 'text-left'}`}>
            <div className={`inline-block px-4 py-2 rounded ${msg.role === 'user' ? 'bg-blue-500 text-white' : 'bg-gray-200 text-black'}`}>
              {msg.content}
            </div>
          </div>
        ))}
      </div>
      <div className="p-2 border-t border-gray-300">
        <input type="text" placeholder="Type a message..." className="w-full border rounded px-2 py-1" />
      </div>
    </div>
  );
}