export default function Sidebar() {
  // Placeholder conversation list
  const conversations = [
    { id: 1, title: 'Chat 1' },
    { id: 2, title: 'Chat 2' },
    { id: 3, title: 'Chat 3' },
  ];

  return (
    <div className="w-64 bg-gray-100 border-r border-gray-300 p-2 overflow-y-auto">
      <h2 className="text-lg font-bold mb-4">Conversations</h2>
      <ul>
        {conversations.map((c) => (
          <li
            key={c.id}
            className="p-2 rounded hover:bg-gray-200 cursor-pointer"
          >
            {c.title}
          </li>
        ))}
      </ul>
    </div>
  );
}
