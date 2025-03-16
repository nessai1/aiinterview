import React, {useState} from "react";
import './QuestionBlock.css';
import MessageEditor from "@/components/internal/chat/editor/MessageEditor.tsx";
import Markdown from "@/components/internal/chat/editor/Markdown.tsx";
import {Question} from "@/lib/interview/interview.ts";
import {
    Accordion,
    AccordionContent,
    AccordionItem,
    AccordionTrigger,
} from "@/components/ui/accordion"



type TProps = {
    question: Question,
    showFeedback: boolean
}



const QuestionBlock: React.FC<TProps> = (props: TProps) => {

    const [answer, setAnswer] = useState<string>(props.question.answer);
    const [isDone, setDone] = useState<boolean>(props.question.done);


    return (
        <div className="questionBlock">
            <div className="question">
                <div className="questionHeader">
                    Вопрос 1:
                </div>
                <div className="questionText">
                    <Markdown>{props.question.text}</Markdown>
                </div>
                <div className="questionMessage">
                    <MessageEditor onAnswer={() => {}} onSkip={() => {}} />
                </div>
                <Accordion className="accordion" type="single" collapsible disabled={props.showFeedback}>
                    <AccordionItem value="item-1">
                        <AccordionTrigger>Обратная связь ИИ</AccordionTrigger>
                        <AccordionContent>
                            <Markdown>{props.question.feedback}</Markdown>
                        </AccordionContent>
                    </AccordionItem>
                </Accordion>
            </div>
        </div>
    );
}

export default QuestionBlock;