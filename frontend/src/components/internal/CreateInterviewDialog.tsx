import React, { FunctionComponent, useState } from 'react';
import { Button } from "@/components/ui/button"
import {
    Dialog,
    DialogContent,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog"

import { MessageCirclePlus } from "lucide-react";
import { Input } from "@/components/ui/input.tsx";
import { Label } from "@/components/ui/label.tsx";
import GradeList from "@/components/internal/GradeList.tsx";

interface GradeItem {
    id: number;
    topic: string;
    grade: string;
}

const CreateInterviewDialog: FunctionComponent = () => {
    const [title, setTitle] = useState("");
    const [time, setTime] = useState("");
    const [grades, setGrades] = useState<GradeItem[]>([{ id: Date.now(), topic: "", grade: "Junior" }]);

    const [errors, setErrors] = useState<{ title?: boolean; time?: boolean; grades?: boolean }>({});

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
            time: time.trim() === "",
            grades: hasEmptyTopic || hasNoTopics
        };

        setErrors(newErrors);

        // Если есть ошибки - не отправляем
        if (Object.values(newErrors).some(error => error)) {
            return;
        }

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
                        className={`mt-1 ${errors.title ? "border-red-500" : ""}`}
                        name="title"
                        placeholder="Собеседование в Google"
                        value={title}
                        onChange={(e) => setTitle(e.target.value)}
                    />
                    {errors.title && <p className="text-red-500 text-sm">Введите название интервью</p>}

                    {/* Тайминг */}
                    <div className="flex items-baseline">
                        <div className="p-2">
                            <Label className="pl-1">Тайминг</Label>
                        </div>
                        <div className="w-24">
                            <Input
                                className={`mt-4 ${errors.time ? "border-red-500" : ""}`}
                                name="date"
                                type="time"
                                value={time}
                                onChange={(e) => setTime(e.target.value)}
                            />
                            {errors.time && <p className="text-red-500 text-sm">Выберите время</p>}
                        </div>
                    </div>

                    {/* Список тем */}
                    <div className="mt-3">
                        <GradeList
                            grades={grades}
                            updateGradeItem={updateGradeItem}
                            addGradeItem={addGradeItem}
                            removeGradeItem={removeGradeItem}
                        />
                        {errors.grades && <p className="text-red-500 text-sm">Добавьте хотя бы одну тему</p>}
                    </div>
                </div>

                <DialogFooter>
                    <Button type="submit" onClick={handleSubmit}>Создать</Button>
                </DialogFooter>

            </DialogContent>
        </Dialog>
    );
};

export default CreateInterviewDialog;
