import { FunctionComponent, useState } from 'react';
import { Button } from "@/components/ui/button"
import {
    Dialog,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"

import {AlarmClock, MessageCirclePlus} from "lucide-react";
import { Input } from "@/components/ui/input.tsx";
import { Label } from "@/components/ui/label.tsx";
import GradeList from "@/components/internal/GradeList.tsx";
import InterviewTimePicker from "@/components/internal/InterviewTimePicker.tsx";
import { Loader2 } from "lucide-react";

interface GradeItem {
    id: number;
    topic: string;
    grade: string;
}

const CreateInterviewDialog: FunctionComponent = () => {
    const [title, setTitle] = useState("");
    const [time, setTime] = useState("");
    const [grades, setGrades] = useState<GradeItem[]>([{ id: Date.now(), topic: "", grade: "Junior" }]);

    const [isLoad, setLoad] = useState(false);
    const [errors, setErrors] = useState<{ title?: boolean; time?: string; grades?: boolean }>({});

    // Функция обновления темы и грейда
    const updateGradeItem = (id: number, field: keyof GradeItem, value: string) => {
        setGrades((prev) =>
            prev.map((item) => (item.id === id ? { ...item, [field]: value } : item))
        );
    };

    // Функция добавления новой темы
    const addGradeItem = () => {
        setGrades([...grades, { id: Date.now(), topic: "", grade: "Junior" }]);
    };

    // Функция удаления темы
    const removeGradeItem = (id: number) => {
        if (grades.length > 1) {
            setGrades(grades.filter((item) => item.id !== id));
        }
    };

    // Функция отправки формы
    const handleSubmit = () => {
        const hasEmptyTopic = grades.some(item => item.topic.trim() === "");
        const hasNoTopics = grades.length === 0;

        const newErrors = {
            title: title.trim() === "",
            time: "",
            grades: hasEmptyTopic || hasNoTopics
        };

        if (time.length <= 0)
        {
            newErrors.time = "Введите время на собеседование (в минутах)";
        }
        else
        {
            const numTime = parseInt(time);
            if (numTime < 5)
            {
                newErrors.time = "Минимальное время собеседования - 5 минут"
            }
            else if ((numTime / grades.length) < 5)
            {
                const gl = grades.length;
                const mt = gl * 5;
                newErrors.time = `На каждую тему должно быть выделено не менее 5 минут. Количество тем - ${gl}, значит минимальный размер собеседования - ${mt} минут`;
            }
        }

        setErrors(newErrors);

        // Если есть ошибки - не отправляем
        if (newErrors.grades || newErrors.title || newErrors.time.length > 0) {
            return;
        }

        setLoad(true);

        // Заглушка запроса (здесь можно сделать реальный API вызов)
        console.log("Отправка данных:", {
            title,
            time,
            grades
        });

        // Закрыть диалог (если требуется)
    };

    return (
        <Dialog>
            <DialogTrigger asChild>
                <Button><MessageCirclePlus /> Новое интервью</Button>
            </DialogTrigger>
            <DialogContent className="sm:max-w-[425px]">

                <DialogHeader>
                    <DialogTitle>Создание интервью</DialogTitle>
                </DialogHeader>

                <div>
                    {/* Название интервью */}
                    <Label className="pl-1">Название интервью</Label>
                    <Input
                        disabled={isLoad}
                        className={`mt-1 ${errors.title ? "border-red-500" : ""}`}
                        name="title"
                        placeholder="Собеседование в Google"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                    />
                    {errors.title && <p className="text-red-500 text-sm">Введите название интервью</p>}

                    {/* Тайминг */}
                    <div className="mt-5 mb-5">
                        <div className="flex items-center">
                            <div className="p-2">
                                <Label className="pl-1 flex"><AlarmClock size={15}/>
                                    <div className={"ml-2"}>Тайминг:</div>
                                </Label>
                            </div>
                            <Input
                                type="number"
                                disabled={isLoad}
                                className={`w-20 h-8 mr-2 ${errors.time ? "border-red-500" : ""}`}
                                onChange={(e) => setTime(e.target.value)}
                            ></Input>
                            <div>
                                минут
                            </div>
                            <div className={"ml-4 text-xs text-zinc-600"}>на одну тему должно быть выделенно минимум 5
                                минут
                            </div>
                        </div>
                        {(errors.time?.length ?? 0) > 0 && <p className="text-red-500 text-sm mt-1">{errors.time}</p>}
                    </div>

                    {/* Список тем */}
                    <div className="mt-3">
                        <GradeList
                            disabled={isLoad}
                            grades={grades}
                            updateGradeItem={updateGradeItem}
                            addGradeItem={addGradeItem}
                            removeGradeItem={removeGradeItem}
                        />
                        {errors.grades && <p className="text-red-500 text-sm">Добавьте хотя бы одну тему</p>}
                    </div>
                </div>

                <DialogFooter>
                    <Button disabled={isLoad}  type="submit" onClick={handleSubmit}>
                        {isLoad ? <Loader2 className="animate-spin" /> : ""}
                        Создать
                    </Button>
                </DialogFooter>

            </DialogContent>
        </Dialog>
    );
};

export default CreateInterviewDialog;
