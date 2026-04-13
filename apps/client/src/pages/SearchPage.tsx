import { useState } from "react";
import { searchNotes } from "../api/client";

export default function SearchPage() {
  const [keyword, setKeyword] = useState("");
  const [count, setCount] = useState(10);
  const [sortType, setSortType] = useState(0);
  const [noteType, setNoteType] = useState(0);
  const [noteTime, setNoteTime] = useState(0);
  const [loading, setLoading] = useState(false);
  const [notes, setNotes] = useState<string[]>([]);
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
    try {
      const data = await searchNotes({
        keyword,
        count,
        sort_type: sortType,
        note_type: noteType,
        note_time: noteTime,
      });
      setNotes(data.notes || []);
    } catch (e: any) {
      setError(e.message);
    } finally {
      setLoading(false);
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
