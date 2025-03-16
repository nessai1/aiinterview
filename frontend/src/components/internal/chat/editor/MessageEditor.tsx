import React, {FunctionComponent, useState, useRef, useEffect} from 'react';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import 'react-tabs/style/react-tabs.css';
import './MessageEditor.css';
import IconButton from "@/components/internal/IconButton.tsx";
import HeadingIcon from "@/assets/icons/editor/heading.svg?react";
import BoldIcon from "@/assets/icons/editor/bold.svg?react";
import ItalicIcon from "@/assets/icons/editor/italic.svg?react";
import QuoteIcon from "@/assets/icons/editor/quote.svg?react";
import CodeIcon from "@/assets/icons/editor/code.svg?react";
import MarkdownIcon from "@/assets/icons/editor/markdown.svg?react";
import Editor from "@/lib/editor.ts";
import Markdown from "./Markdown.tsx";
import {Button} from "@/components/ui/button.tsx";
import {Loader2} from "lucide-react";


interface TProps {
    onSkip: (onComplete: () => void) => void,
    onAnswer: (onComplete: () => void, answer: string) => void
}

const MessageEditor: FunctionComponent<TProps> = (props: TProps) => {

    const [isEditing, setEditing] = useState<boolean>(true);

    const textareaRef = useRef<HTMLTextAreaElement | null>(null);

    const [previewContent, setPreviewContent] = useState<string>('Загрузка...');
    const [previousContent, setPreviousContent] = useState<string>('');

    const [isLoading, setLoading] = useState<boolean>(false);

    const [editor, setEditor] = useState<Editor | null>(null);

    useEffect(() => {
        if (textareaRef.current instanceof HTMLTextAreaElement && editor === null) {
            setEditor(new Editor(textareaRef.current));
        }
    }, [textareaRef]);

    const onComplete = () => {

    }

    return (
        <div>
            <Tabs className="editor">
                <TabList className="tab-list">
                    <div>
                        <Tab onClick={() => setEditing(true)}>Редактирование</Tab>
                        <Tab onClick={() => {
                            setEditing(false);

                            const content = editor?.getContent() ?? '';
                            if (content !== previousContent) {
                                setPreviewContent('Загрузка...');
                                window.network.previewMessage(editor?.getContent() ?? '')
                                    .then((val: string) => {
                                        setPreviewContent(val);
                                        setPreviousContent(content);
                                    })
                                    .catch((e) => {
                                        console.log(e);
                                    });
                            }

                        }}>Предпросмотр
                        </Tab>
                    </div>
                    <div className="action-list" style={{display: isEditing ? 'flex' : 'none'}}>
                        <IconButton onClick={() => {
                            editor?.heading();
                        }} tooltip={"Заголовок"}>
                            <HeadingIcon/>
                        </IconButton>
                        <IconButton onClick={() => {
                            editor?.bold();
                        }} tooltip={"Bold"}>
                            <BoldIcon/>
                        </IconButton>
                        <IconButton onClick={() => {
                            editor?.italic();
                        }} tooltip={"Italic"}>
                            <ItalicIcon/>
                        </IconButton>
                        <IconButton onClick={() => {
                            editor?.quotes();
                        }} tooltip={"Цитата"}>
                            <QuoteIcon/>
                        </IconButton>
                        <IconButton onClick={() => {
                            editor?.code();
                        }} tooltip={"Код"}>
                            <CodeIcon/>
                        </IconButton>
                    </div>
                </TabList>
                <TabPanel></TabPanel>
                <div className="editor-active" style={{display: isEditing ? 'block' : 'none'}}>
                    <div className="comment-input-wrapper">
                    <textarea name="comment" id="" className="comment-input" ref={textareaRef}
                              placeholder="Напишите ответ на вопрос" disabled={isLoading}>
                    </textarea>
                    </div>
                    <div className="editor-bottom"><MarkdownIcon/> Есть поддержка Markdown</div>
                </div>
                <TabPanel>
                    <Markdown>
                        {previewContent}
                    </Markdown>
                </TabPanel>
            </Tabs>
            <div className="flex justify-end gap-2 mt-2 items-center">
                {isLoading ? <Loader2 className="animate-spin" /> : ""}
                <Button onClick={() => {
                    const content = editor?.getContent() ?? '';
                    setLoading(true);
                    props.onAnswer(onComplete, content);
                }} className="button-ok" disabled={isLoading}>Ответить</Button>
                <Button onClick={() => {
                    setLoading(true);
                    props.onSkip(onComplete);
                }} className="button-decline" disabled={isLoading}>Пропустить вопрос</Button>
            </div>
        </div>
    );
};

export default MessageEditor;
