"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { X } from "lucide-react"; // Иконка крестика

const GRADE_OPTIONS = ["Junior", "Middle", "Senior"];

export default function GradeList() {
    const [items, setItems] = useState([{ id: Date.now(), topic: "", grade: "Junior" }]);

    // Добавление нового пункта
    const addItem = () => {
        setItems([...items, { id: Date.now(), topic: "", grade: "Junior" }]);
    };

    // Удаление пункта
    const removeItem = (id) => {
        if (items.length > 1) {
            setItems(items.filter((item) => item.id !== id));
        }
    };

    // Обновление значений инпутов
    const updateItem = (id, field, value) => {
        setItems((prev) =>
            prev.map((item) => (item.id === id ? { ...item, [field]: value } : item))
        );
    };

    return (
        <div className="space-y-4 p-4 border rounded-lg shadow-sm">
            <h2 className="text-lg font-semibold">Список тем</h2>
            {items.map((item) => (
                <div key={item.id} className="flex gap-4 items-center">
                    <Input
                        type="text"
                        placeholder="Введите тему"
                        value={item.topic}
                        onChange={(e) => updateItem(item.id, "topic", e.target.value)}
                        className="flex-1"
                    />
                    <Select
                        value={item.grade}
                        onValueChange={(value) => updateItem(item.id, "grade", value)}
                    >
                        <SelectTrigger className="w-36">
                            <SelectValue placeholder="Выберите грейд" />
                        </SelectTrigger>
                        <SelectContent>
                            {GRADE_OPTIONS.map((grade) => (
                                <SelectItem key={grade} value={grade}>
                                    {grade}
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                    {items.length > 1 && (
                        <Button variant="ghost" size="icon" onClick={() => removeItem(item.id)}>
                            <X className="w-5 h-5 text-red-500" />
                        </Button>
                    )}
                </div>
            ))}
            <Button onClick={addItem} className="w-full mt-2">
                + Добавить пункт
            </Button>
        </div>
    );
}
