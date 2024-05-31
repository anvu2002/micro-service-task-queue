import React, { useState, useRef } from 'react';
import axios from 'axios';
import './App.css';

const App = () => {
  const [prompt, setPrompt] = useState('');
  const [progress, setProgress] = useState(0);
  const [statusList, setStatusList] = useState([]);
  const [imageTaskID, setImageTaskID] = useState(null);
  const [sentenceTaskID, setSentenceTaskID] = useState(null);

  const intervalRef = useRef(null);

  const URL_GO = 'http://localhost:9888';
  const URL_GO_IMAGE_STATUS = `${URL_GO}/image_status?task_id=`;
  const URL_GO_KEYWORD_STATUS = `${URL_GO}/keyword_status?task_id=`;
  const URL_GO_TTS_STATUS = `${URL_GO}/tts_status?task_id=`;

  const handleGetVideo = () => {
    const file_name = encodeURIComponent(prompt);

    axios.post(`${URL_GO}/process_doc?file_name=${file_name}`)
      .then(response => {
        const { go_task_id } = response.data;
        setStatusList([go_task_id]);
        console.log("process_doc task_id:", go_task_id); // Log the task_id
        intervalRef.current = setInterval(refresh, 1000);
      })
      .catch(error => {
        console.error("Error processing document:", error);
      });
  };

  const refresh = () => {
    console.log("Refreshing statuses..."); // Log refresh action
    const promises = statusList.map(taskId => {
      console.log(`Checking status for task: ${taskId}`); // Log each status check
      return axios.get(`${URL_GO_KEYWORD_STATUS}${taskId}`)
        .catch(error => {
          console.error(`Error fetching keyword status for task ${taskId}:`, error);
          return { data: { status: 'FAILED' } }; // Return a failed status if there is an error
        });
    });

    Promise.all(promises)
      .then(responses => {
        let n_success = 0;
        responses.forEach((response) => {
          const { status, keywords } = response.data;
          console.log(`Keyword status for task: ${status}`); // Log the status
          if (status === 'SUCCESS' || status === 'FAILED') {
            n_success++;
          }
          if (status === 'SUCCESS') {
            continueTasks(keywords);
          }
        });

        setProgress(prev => (prev < 100 ? prev + Math.floor(Math.random() * 10) + 1 : 100));

        if (n_success === statusList.length) {
          clearInterval(intervalRef.current);
          setProgress(0);
        }
      })
      .catch(error => {
        console.error("Error refreshing status:", error);
      });
  };

  const continueTasks = (keywords) => {
    console.log("Continuing tasks with keywords:", keywords); // Log keywords

    axios.get(`${URL_GO}/get_images?query=${encodeURIComponent(keywords['keywords'])}`)
      .then(response => {
        const { go_task_id } = response.data;
        setImageTaskID(go_task_id);
        console.log("get_images task_id:", go_task_id); // Log the image task_id
      })
      .catch(error => {
        console.error("Error fetching images:", error);
      });

    axios.get(`${URL_GO}/tts?text=${encodeURIComponent(keywords['sentences'])}&save_path=./`)
      .then(response => {
        const { go_task_id } = response.data;
        setSentenceTaskID(go_task_id);
        console.log("tts task_id:", go_task_id); // Log the TTS task_id
      })
      .catch(error => {
        console.error("Error generating TTS:", error);
      });

    intervalRef.current = setInterval(checkTaskStatus, 1000);
  };

  const checkTaskStatus = () => {
    console.log("Checking image and TTS task statuses..."); // Log check action
    const promises = [
      axios.get(`${URL_GO_IMAGE_STATUS}${imageTaskID}`)
        .catch(error => {
          console.error(`Error fetching image status for task ${imageTaskID}:`, error);
          return { data: { status: 'FAILED' } }; // Return a failed status if there is an error
        }),
      axios.get(`${URL_GO_TTS_STATUS}${sentenceTaskID}`)
        .catch(error => {
          console.error(`Error fetching TTS status for task ${sentenceTaskID}:`, error);
          return { data: { status: 'FAILED' } }; // Return a failed status if there is an error
        })
    ];

    Promise.all(promises)
      .then(responses => {
        const imageStatus = responses[0].data.status;
        const sentenceStatus = responses[1].data.status;
        console.log(`Image status: ${imageStatus}, TTS status: ${sentenceStatus}`); // Log the statuses

        if (imageStatus === 'SUCCESS' && sentenceStatus === 'SUCCESS') {
          clearInterval(intervalRef.current);
          ffmpegTask();
        }
      })
      .catch(error => {
        console.error("Error checking task status:", error);
      });
  };

  const ffmpegTask = () => {
    axios.get(`${URL_GO}/get_ffmpeg?audio_path=./audio&image_path=./images`)
      .then(response => {
        console.log("[Final Task] ffmpeg Completed", response.data);
      })
      .catch(error => {
        console.error("Error in ffmpeg task:", error);
      });
  };

  return (
    <div>
      <input
        type="text"
        value={prompt}
        onChange={e => setPrompt(e.target.value)}
        placeholder="Enter file prompt"
      />
      <button onClick={handleGetVideo}>Get Video</button>
      <div id="progress-bar" style={{ width: `${progress}%`, backgroundColor: 'blue' }}>
        {progress}%
      </div>
    </div>
  );
};

export default App;
