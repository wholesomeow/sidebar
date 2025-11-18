export default function ChatPanel() {
  const messages = [
    { id: 1, role: 'assistant', content: 'Hello! How can I help you?' },
    { id: 2, role: 'user', content: 'I want to build a React UI.' },
  ];

  return (
    <div className="flex-1 flex flex-col bg-white">
      <div className="flex-1 p-4 overflow-y-auto">
        {messages.map((msg) => (
          <div
            key={msg.id}
            className={`mb-4 ${
              msg.role === 'user' ? 'text-right' : 'text-left'
            }`}
          >
            <div
              className={`inline-block px-4 py-2 rounded ${
                msg.role === 'user'
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-200 text-black'
              }`}
            >
              {msg.content}
            </div>
          </div>
        ))}
      </div>
      <div className="p-2 border-t border-gray-300">
        <input
          type="text"
          placeholder="Type a message..."
          className="w-full border rounded px-2 py-1"
        />
      </div>
    </div>
  );
}
