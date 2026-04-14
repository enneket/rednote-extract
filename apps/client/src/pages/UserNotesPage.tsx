import { useState } from "react";
import { collectUserNotes, exportNotes, NoteData } from "../api/client";

export default function UserNotesPage() {
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [notes, setNotes] = useState<NoteData[]>([]);
  const [partial, setPartial] = useState(false);
  const [errors, setErrors] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);

  async function handleCollect() {
    if (!url.trim()) {
      setError("请输入用户主页 URL");
      return;
    }
    setLoading(true);
    setError(null);
    setNotes([]);
    setErrors([]);
    try {
      const data = await collectUserNotes(url);
      setNotes(data.notes);
      setPartial(data.partial);
      setErrors(data.errors);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleExport() {
    try {
      await exportNotes(notes, "user_notes.xlsx");
    } catch (e: any) {
      setError(e.message);
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
        {notes.length > 0 && (
          <button className="btn btn-secondary" onClick={handleExport}>
            导出Excel
          </button>
        )}
      </div>

      {error && <div className="error" style={{ marginTop: 12 }}>{error}</div>}

      {notes.length > 0 && (
        <div className="result-box">
          <div className="title">
            采集到 {notes.length} 条笔记
            {partial && <span style={{ color: "#e53935", fontSize: 13 }}>（部分失败）</span>}
          </div>
          <div className="note-list">
            {notes.map((note, i) => (
              <div key={i} className="note-item">
                <p><strong>{note.title || "(无标题)"}</strong> - {note.author}</p>
                <p>点赞:{note.liked} 收藏:{note.collected} 评论:{note.commented}</p>
                <a href={note.url} target="_blank" rel="noopener noreferrer" className="url">{note.url}</a>
              </div>
            ))}
          </div>
          {errors.length > 0 && (
            <div style={{ marginTop: 12, fontSize: 12, color: "#e53935" }}>
              失败: {errors.length} 条
            </div>
          )}
        </div>
      )}
    </div>
  );
}
