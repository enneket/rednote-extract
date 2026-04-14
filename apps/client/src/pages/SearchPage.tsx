import { useState } from "react";
import { searchNotes, exportNotes, NoteData } from "../api/client";

export default function SearchPage() {
  const [keyword, setKeyword] = useState("");
  const [count, setCount] = useState(10);
  const [sortType, setSortType] = useState(0);
  const [noteType, setNoteType] = useState(0);
  const [noteTime, setNoteTime] = useState(0);
  const [loading, setLoading] = useState(false);
  const [notes, setNotes] = useState<NoteData[]>([]);
  const [partial, setPartial] = useState(false);
  const [errors, setErrors] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);

  const sortOptions = [
    { value: 0, label: "综合排序" },
    { value: 1, label: "最新" },
    { value: 2, label: "最多点赞" },
    { value: 3, label: "最多评论" },
    { value: 4, label: "最多收藏" },
  ];

  const noteTypeOptions = [
    { value: 0, label: "不限" },
    { value: 1, label: "视频笔记" },
    { value: 2, label: "普通笔记" },
  ];

  const noteTimeOptions = [
    { value: 0, label: "不限" },
    { value: 1, label: "一天内" },
    { value: 2, label: "一周内" },
    { value: 3, label: "半年内" },
  ];

  async function handleSearch() {
    if (!keyword.trim()) {
      setError("请输入搜索关键词");
      return;
    }
    setLoading(true);
    setError(null);
    setNotes([]);
    setErrors([]);
    try {
      const data = await searchNotes({
        keyword,
        count,
        sort_type: sortType,
        note_type: noteType,
        note_time: noteTime,
      });
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
      await exportNotes(notes, `${keyword}_notes.xlsx`);
    } catch (e: any) {
      setError(e.message);
    }
  }

  return (
    <div>
      <div className="page-title">搜索采集</div>

      <div className="form-group">
        <label>搜索关键词</label>
        <input
          type="text"
          placeholder="输入关键词..."
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
        />
      </div>

      <div className="form-row">
        <div className="form-group">
          <label>采集数量</label>
          <input
            type="number"
            min={1}
            max={100}
            value={count}
            onChange={(e) => setCount(Number(e.target.value))}
          />
        </div>
        <div className="form-group">
          <label>排序方式</label>
          <select value={sortType} onChange={(e) => setSortType(Number(e.target.value))}>
            {sortOptions.map((o) => (
              <option key={o.value} value={o.value}>{o.label}</option>
            ))}
          </select>
        </div>
      </div>

      <div className="form-row">
        <div className="form-group">
          <label>笔记类型</label>
          <select value={noteType} onChange={(e) => setNoteType(Number(e.target.value))}>
            {noteTypeOptions.map((o) => (
              <option key={o.value} value={o.value}>{o.label}</option>
            ))}
          </select>
        </div>
        <div className="form-group">
          <label>笔记时间</label>
          <select value={noteTime} onChange={(e) => setNoteTime(Number(e.target.value))}>
            {noteTimeOptions.map((o) => (
              <option key={o.value} value={o.value}>{o.label}</option>
            ))}
          </select>
        </div>
      </div>

      <div className="btn-group">
        <button className="btn btn-primary" onClick={handleSearch} disabled={loading}>
          {loading ? "搜索中..." : "搜索并采集"}
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
