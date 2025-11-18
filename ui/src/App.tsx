import Sidebar from './components/Sidebar';
import ChatPanel from './components/ChatPanel';
import RightPanel from './components/RightPanel';

export default function App() {
    return (
        <div className='h-screen flex overflow-hidden'>
            {/* Left Sidebar */}
            <Sidebar />

            {/*Main Chat */}
            <ChatPanel />

            {/* Right Sidebar */}
            <RightPanel />
        </div>
    );
}