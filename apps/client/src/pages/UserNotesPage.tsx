import { useState } from "react";
import { collectUserNotes } from "../api/client";

export default function UserNotesPage() {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [notes, setNotes] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);

  async function handleCollect() {
    if (!url.trim()) {
      setError("请输入用户主页 URL");
      return;
    }
    setLoading(true);
    setError(null);
    setNotes([]);
    try {
      const data = await collectUserNotes(url);
      setNotes(data.notes || []);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div>
      <div className="page-title">用户笔记采集</div>

      <div className="form-group">
        <label>用户主页 URL</label>
        <input
          type="text"
          placeholder="https://www.xiaohongshu.com/user/profile/..."
          value={url}
          onChange={(e) => setUrl(e.target.value)}
        />
      </div>

      <div className="btn-group">
        <button className="btn btn-primary" onClick={handleCollect} disabled={loading}>
          {loading ? "采集中..." : "采集全部"}
        </button>
      </div>

      {error && <div className="error" style={{ marginTop: 12 }}>{error}</div>}

      {notes.length > 0 && (
        <div className="result-box">
          <div className="title">采集到 {notes.length} 条笔记</div>
          <div className="note-list">
            {notes.map((noteUrl, i) => (
              <div key={i} className="note-item">
                <span className="url">{noteUrl}</span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
