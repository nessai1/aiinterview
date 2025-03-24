import React, {useEffect, useRef} from "react";
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
import {useToast} from "@/hooks/use-toast.ts";
import {AxiosError} from "axios";



type TProps = {
    question: Question,
    onNextSection: (q: Question) => void,
    onNextQuestion: (q: Question) => void,
}



const QuestionBlock: React.FC<TProps> = (props: TProps) => {

    const question = props.question;

    const { toast } = useToast();

    const formRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        if (!props.question.done && formRef.current instanceof HTMLDivElement) {
            formRef.current.scrollIntoView({behavior: "smooth"});
        }
    }, []);

    const answerFunc = (answer: string, onComplete: () => void) => {
        network.answerQuestion(question.uuid, answer).then((response) => {
            if (response.updatedQuestion == null)
            {
                console.log("not loaded question", response);
                return;
            }

            if (response.nextSection)
            {
                props.onNextSection(response.updatedQuestion);
            }
            else
            {
                props.onNextQuestion(response.updatedQuestion);
            }
        }).catch((err: AxiosError) => {
            onComplete();
            toast({
                title: 'Упс! Ответ не отправлен ;(',
                description: `Ошибка сети: [${err.code}] ${err.message}`,
                variant: "destructive",
            });
        });
    }

    return (
        <div className="questionBlock">
            <div className="question">
                <div className="questionHeader">
                    Вопрос:
                </div>
                <div className="questionText">
                    <Markdown>{props.question.question}</Markdown>
                </div>
                <div className="questionHeader">
                    Ответ:
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
                        <div className="questionMessage" ref={formRef}>
                            <MessageEditor
                                onAnswer={(onComplete, answer) => {
                                   answerFunc(answer, onComplete);
                                }}
                                onSkip={(onComplete) => {
                                    answerFunc("", onComplete);
                                }} />
                        </div>
                    )
                }

            </div>
        </div>
    );
}

export default QuestionBlock;