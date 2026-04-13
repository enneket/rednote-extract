import { useState } from "react";
import CookiePage from "./pages/CookiePage";
import SingleNotePage from "./pages/SingleNotePage";
import UserNotesPage from "./pages/UserNotesPage";
import SearchPage from "./pages/SearchPage";

type Tab = "cookie" | "single" | "user" | "search";

export default function App() {
  const [activeTab, setActiveTab] = useState<Tab>("cookie");

  return (
    <div className="app">
      <header className="header">
        <h1>XHS 数据采集</h1>
      </header>

      <nav className="nav">
        <button
          className={activeTab === "cookie" ? "active" : ""}
          onClick={() => setActiveTab("cookie")}
        >
          Cookie
        </button>
        <button
          className={activeTab === "single" ? "active" : ""}
          onClick={() => setActiveTab("single")}
        >
          单条笔记
        </button>
        <button
          className={activeTab === "user" ? "active" : ""}
          onClick={() => setActiveTab("user")}
        >
          用户笔记
        </button>
        <button
          className={activeTab === "search" ? "active" : ""}
          onClick={() => setActiveTab("search")}
        >
          搜索采集
        </button>
      </nav>

      <main className="content">
        {activeTab === "cookie" && <CookiePage />}
        {activeTab === "single" && <SingleNotePage />}
        {activeTab === "user" && <UserNotesPage />}
        {activeTab === "search" && <SearchPage />}
      </main>
    </div>
  );
}
