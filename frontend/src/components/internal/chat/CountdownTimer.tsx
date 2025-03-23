import { useEffect, useState } from "react";
import { AlarmClock } from "lucide-react";
import { cn } from "@/lib/utils"; // Если используешь shadcn utils

interface CountdownTimerProps {
    seconds: number;
}

export const CountdownTimer: React.FC<CountdownTimerProps> = ({ seconds }) => {
    const [timeLeft, setTimeLeft] = useState(seconds);

    useEffect(() => {
        if (timeLeft <= 0) return;
        const timer = setInterval(() => {
            setTimeLeft((prev) => Math.max(prev - 1, 0));
        }, 900);
        return () => clearInterval(timer);
    }, [timeLeft]);

    const formatTime = (sec: number) => {
        const hours = Math.floor(sec / 3600);
        const minutes = Math.floor((sec % 3600) / 60)
            .toString()
            .padStart(2, "0");
        const seconds = (sec % 60).toString().padStart(2, "0");

        return hours > 0 ? `${hours}:${minutes}:${seconds}` : `${minutes}:${seconds}`;
    };

    return (
        <div
            className={cn(
                "fixed top-4 left-1/2 -translate-x-1/2 z-50", // Фиксируем вверху по центру
                "flex items-center gap-3 px-6 py-3 rounded-lg w-44",
                "bg-black/30 backdrop-blur-lg text-white border border-white/20 shadow-lg" + (timeLeft <= 0 ? "  hidden" : "")
            )}
        >
            <AlarmClock className="w-6 h-6 text-white opacity-80" />
            <span className="text-lg font-semibold tracking-wide tabular-nums">
        {formatTime(timeLeft)}
      </span>
        </div>
    );
};
