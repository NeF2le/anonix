function ttlToHuman(proto) {
    if (!proto) return "";
    const seconds = Number(proto.seconds || 0);
    const nanos = Number(proto.nanos || 0);

    let totalSeconds = seconds + Math.floor(nanos / 1e9);
    if (totalSeconds === 0 && nanos > 0 && nanos < 1e7) {
        totalSeconds = nanos;
    }

    if (totalSeconds <= 0) return "";

    const hours = Math.floor(totalSeconds / 3600);
    const minutes = Math.floor((totalSeconds % 3600) / 60);
    const secs = totalSeconds % 60;

    if (hours > 0) return `${hours}h ${minutes}m`;
    if (minutes > 0) return `${minutes}m ${secs}s`;
    return `${secs}s`;
}

function timestampToDate(ts) {
    if (!ts) return null;
    const sec = Number(ts.seconds || 0);
    const nanos = Number(ts.nanos || 0);
    if (!Number.isFinite(sec) || !Number.isFinite(nanos)) return null;
    const ms = sec * 1000 + Math.floor(nanos / 1e6);
    return new Date(ms);
}