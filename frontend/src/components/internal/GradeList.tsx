import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { X } from "lucide-react"; // Иконка крестика

const GRADE_OPTIONS = ["Junior", "Middle", "Senior"];

export default function GradeList({ grades, updateGradeItem, addGradeItem, removeGradeItem }) {
    return (
        <div className="space-y-2">
            {grades.map((item) => (
                <div key={item.id} className="flex gap-4 items-center">
                    <Input
                        type="text"
                        placeholder="Введите тему"
                        value={item.topic}
                        onChange={(e) => updateGradeItem(item.id, "topic", e.target.value)}
                        className={`flex-1 ${item.topic.trim() === "" ? "border-red-500" : ""}`}
                    />
                    <Select
                        value={item.grade}
                        onValueChange={(value) => updateGradeItem(item.id, "grade", value)}
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
                    {grades.length > 1 && (
                        <Button variant="ghost" size="icon" onClick={() => removeGradeItem(item.id)}>
                            <X className="w-5 h-5 text-red-500" />
                        </Button>
                    )}
                </div>
            ))}
            <Button onClick={addGradeItem} className="w-full mt-2">
                + Добавить тему
            </Button>
        </div>
    );
}
