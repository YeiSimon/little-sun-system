```mermaid

flowchart LR
    %% 設定圖表為從左到右(LR)而非上到下(TD)，並優化連線
    %% 增加初始化設定控制圖表大小
    %%{init: {'flowchart': {'nodeSpacing': 50, 'rankSpacing': 70, 'curve': 'linear'}}}%%
    User([使用者]) <--> Frontend[前端界面]
    Frontend <--> APIGateway[API Gateway\n+ Rate Limit]
    APIGateway --> MaliciousCheck[惡意使用者語意辨識系統]
    MaliciousCheck --> |安全檢查通過| AIAssistant[AI智慧助理系統語意相似度比對]
    QdrantDB[(QA向量資料庫)] --> |資料提供| AIAssistant
    BigData[(大數據中心之資料)] --> |資料提供| AIAssistant
    KnowledgeBase[(知識庫)] --> |資料提供| AIAssistant
    MaliciousCheck --> |檢測到惡意| Reject([拒絕請求])
    AIAssistant --> Cache[Cache 緩衝系統]
    Cache --> |緩存命中| Frontend
    Cache --> |緩存未命中| Ollama[Ollama 大語言模型]
    Ollama --> Cache
    Cache --> APIGateway
    
    %% 定義樣式
    classDef system fill:#f9f,stroke:#333,stroke-width:1.5px;
    classDef server fill:#bbf,stroke:#333,stroke-width:1.5px;
    classDef user fill:#dfd,stroke:#333,stroke-width:1.5px;
    classDef security fill:#fdd,stroke:#333,stroke-width:1.5px;
    classDef database fill:#ffd,stroke:#333,stroke-width:1.5px;
    
    %% 設定各節點樣式大小
    style User width:120px,height:50px;
    style Frontend width:120px,height:50px;
    style APIGateway width:140px,height:60px;
    style MaliciousCheck width:140px,height:60px;
    style AIAssistant width:160px,height:60px;
    style Cache width:140px,height:50px;
    style Ollama width:160px,height:50px;
    style QdrantDB width:130px,height:60px;
    style BigData width:130px,height:60px;
    style KnowledgeBase width:130px,height:60px;
    
    class User user;
    class Frontend,APIGateway server;
    class MaliciousCheck security;
    class AIAssistant,Cache,Ollama system;
    class QdrantDB,BigData,KnowledgeBase database;
    
    ```