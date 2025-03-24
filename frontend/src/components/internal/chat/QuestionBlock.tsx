import React from "react";
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
    onAnswer: () => void,
    showFeedback: boolean
}



const QuestionBlock: React.FC<TProps> = (props: TProps) => {
    return (
        <div className="questionBlock">
            <div className="question">
                <div className="questionHeader">
                    Вопрос 1:
                </div>
                <div className="questionText">
                    <Markdown>{props.question.question}</Markdown>
                </div>
                {
                    props.question.done ? (
                        <>
                            <div className="questionMessage">
                                {
                                    props.question.answer === "" ? (
                                        <div className="noAnswer">
                                            Ответ не предоставлен
                                        </div>
                                    ) : <Markdown>{props.question.answer}</Markdown>
                                }
                            </div>
                            <Accordion className="accordion" type="single" collapsible>
                                <AccordionItem value="item-1">
                                    <AccordionTrigger>Обратная связь ИИ</AccordionTrigger>
                                    <AccordionContent>
                                        <Markdown>{props.question.feedback}</Markdown>
                                    </AccordionContent>
                                </AccordionItem>
                            </Accordion>
                        </>
                    ) : (
                        <div className="questionMessage">
                            <MessageEditor onAnswer={() => {}} onSkip={() => {}} />
                        </div>
                    )
                }

            </div>
        </div>
    );
}

export default QuestionBlock;