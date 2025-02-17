import TimePicker from "react-time-picker";
import "react-time-picker/dist/TimePicker.css";
import "./InterviewTimePicker.css";
import {AlarmClock} from "lucide-react";


type TProps = {
    time: string,
    setTime: (time: string|null) => void,
    hasError: boolean,
    disabled: boolean,
};

export default function InterviewTimePicker({ time, setTime, hasError, disabled }: TProps) {
    return (
        <div className="w-sm">
            <div className={"flex items-center"}>
                <TimePicker
                    disabled={disabled}
                    onChange={setTime}  // Обновляет состояние времени
                    value={time}         // Текущее значение времени
                    disableClock={true}  // Убирает выбор через аналоговые часы
                    format="HH:mm"       // Принудительно ставит 24-часовой формат
                    clearIcon={null}     // Убирает кнопку очистки

                    className={`border rounded-md p-1  text-sm mr-1  ${
                        hasError ? "border-red-500" : "border-gray-300"
                    }`}
                />
                <AlarmClock size={15}/>
            </div>
            {hasError && <p className="text-red-500 text-xs mt-1">Выберите время</p>}
        </div>
    );
}
