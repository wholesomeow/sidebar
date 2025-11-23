import { useState } from 'react';
import ConversationPanel from './components/ConversationPanel';
import ChatPanel from './components/ChatPanel';
import RightPanel from './components/RightPanel';

export default function App() {
    const [selectedConversationID, setSelectedConversationID] = useState<string | null>(null);

    return (
        <div className='h-screen flex overflow-hidden'>
            {/* Left Sidebar */}
            <ConversationPanel onSelectConversation={setSelectedConversationID} />

            {/*Main Chat */}
            <ChatPanel conversationID={selectedConversationID} />

            {/* Right Sidebar */}
            <RightPanel />
        </div>
    );
}