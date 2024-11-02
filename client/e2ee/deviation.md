Expectation:
// NOTES ////
/*
Messages become only Signal messages when both users have acknowledged sessions in both directions. This means that each user has:

Established a session with the other user for sending messages.
Acknowledged the session by receiving and processing a message from the other user.
Detailed Message Flow Between User1 and User2
Let's break down the message exchange step by step:

Initial State
No sessions are established between user1 and user2.
Both users need to establish sessions for sending messages to each other.
Sequence of Messages
User1 Sends First PreKey Message to User2

Message: PreKey Message
From: user1
To: user2
Purpose: User1 initiates a session for sending messages to user2.
User1's Session State: Unacknowledged session to user2.
User2's Session State: No session established yet.
User1 Sends Second PreKey Message to User2

Message: PreKey Message
From: user1
To: user2
Purpose: Since user1 hasn't received a response, the session remains unacknowledged, so user1 continues sending PreKey messages.
User1's Session State: Unacknowledged session to user2.
User2 Receives Messages from User1 and Processes Them

Action: User2 receives and processes user1's PreKey messages.
User2's Session State: Acknowledged session from user1 to user2 (for decrypting messages from user1).
User2 Sends First PreKey Message to User1

Message: PreKey Message
From: user2
To: user1
Purpose: User2 initiates a session for sending messages to user1.
User2's Session State: Unacknowledged session to user1.
User1's Session State: Acknowledged session to user2, unacknowledged session from user2.
User2 Sends Second PreKey Message to User1

Message: PreKey Message
From: user2
To: user1
Purpose: Since user2 hasn't received a response from user1, the session remains unacknowledged.
User2's Session State: Unacknowledged session to user1.
User1 Receives Messages from User2 and Processes Them

Action: User1 receives and processes user2's PreKey messages.
User1's Session State: Acknowledged session from user2 to user1 (for decrypting messages from user2).
User1 Sends Signal Message to User2

Message: Signal Message
From: user1
To: user2
Purpose: Now that user1 has processed user2's messages, the session is acknowledged in both directions from user1's perspective.
User1's Session State: Acknowledged session to and from user2.
User2's Session State: Unacknowledged session to user1.
User1 Sends Another Signal Message to User2

Message: Signal Message
From: user1
To: user2
Purpose: Continued communication using Signal messages.
User1's Session State: Acknowledged in both directions.
User2 Receives Signal Messages from User1 and Processes Them

Action: User2 receives and processes user1's Signal messages.
User2's Session State: Acknowledged session from user1 to user2, and now acknowledged session to user1 after processing the message.
User2 Sends Signal Message to User1

Message: Signal Message
From: user2
To: user1
Purpose: With the session now acknowledged in both directions, user2 switches to Signal messages.
User2's Session State: Acknowledged in both directions.
User2 Sends Another Signal Message to User1

Message: Signal Message
From: user2
To: user1
Purpose: Continued communication using Signal messages.
Both Users' Session States: Fully acknowledged in both directions.
Summary
PreKey Messages Continue: Each user continues to send PreKey messages until they have processed a message from the other user in the respective direction.
Switch to Signal Messages: Once a user processes a message from the other user, they can start sending Signal messages in that direction.
Clarification of Your Sequence
Your sequence:

user1 prekey -> user2
user1 prekey -> user2
user2 prekey -> user1
user2 prekey -> user1
user1 signal -> user2
user1 signal -> user2
user2 signal -> user1
user2 signal -> user1
*/




A DEVIATION:
Certainly! Let's delve into your test code and understand why the first message **User1** receives from **User2** is being decrypted as a **Signal message** instead of a **PreKey message**.

### **Summary of Your Observation**

- **Issue**: When **User1** tries to decrypt the first message from **User2** as a **PreKey message**, it fails.
- **Observation**: Decrypting it as a **Signal message** succeeds.
- **Question**: "Looks like the first message **User1** receives from **User2** is a Signal message... is this wrong?"

### **Understanding the Expected Behavior**

Based on the Signal Protocol, the expected message flow when two users establish communication is:

1. **User1** sends a **PreKey message** to **User2**.
2. **User2** processes the **PreKey message**, establishes a session from **User1** to **User2**, and then sends a **PreKey message** back to **User1** to establish the session from **User2** to **User1**.
3. Both users continue to exchange **Signal messages** after the sessions are acknowledged.

### **Analyzing Your Test Code**

Let's look at the specific part of your test code:

```go
// User1 receives and decrypts the PreKey message from User2
select {
case msg2 := <-client1.ChatClient.MessageChannel:
    t.Logf("User1 received message from User2: %v", msg2)
    // decryptedMsg2, err := client1.ChatClient.DecryptMessage(context.Background(), msg2, "prekey")
    decryptedMsg2, err := client1.ChatClient.DecryptMessage(context.Background(), msg2, "signal")
    if err != nil {
        t.Fatalf("User1 failed to decrypt message: %v", err)
    }
    t.Logf("User1 decrypted message: %s", decryptedMsg2)
    if decryptedMsg2 != messageFromUser2 {
        t.Fatalf("User1 decrypted message does not match original message: %s", decryptedMsg2)
    }
```

- **Commented Line**: You commented out the line where you attempt to decrypt the message as a **PreKey message**.
- **Active Line**: Instead, you're decrypting it as a **Signal message**, and it succeeds.

### **Why Is This Happening?**

#### **1. Session State Management**

- **User2's Session State**: After **User2** receives and processes **User1**'s **PreKey message**, it establishes a session **from User1 to User2**.
- **Session Acknowledgment**: Some implementations may consider the session as **fully established** after processing the initial **PreKey message** from the other user.
- **Sending Messages**: **User2**'s client may now believe that it has an **established and acknowledged session** with **User1**, and therefore sends a **Signal message** instead of a **PreKey message**.

#### **2. Implementation Details**

- **Optimistic Session Establishment**: Your implementation might optimistically assume that since it has processed a **PreKey message** from **User1**, the session is established in both directions.
- **Sending Signal Messages Early**: As a result, **User2** sends a **Signal message** back to **User1** immediately, skipping the step of sending a **PreKey message**.

#### **3. Decryption Fails as PreKey Message**

- **Mismatch in Message Type**: When **User1** tries to decrypt **User2**'s message as a **PreKey message**, it fails because the message is actually a **Signal message**.
- **Successful Decryption as Signal Message**: Decrypting it as a **Signal message** succeeds because that's the actual message type.

### **Is This Wrong?**

#### **According to the Signal Protocol Specification**

- **Sessions Are Directional**: Each user must establish a session for sending messages to the other user.
- **First Message Should Be PreKey Message**: The first message from **User2** to **User1** should be a **PreKey message** if no session has been established in that direction.

#### **Possible Deviations in Implementations**

- **Optimizations**: Some implementations might optimize the protocol flow by considering the session established in both directions after receiving a **PreKey message**.
- **Security Implications**: While this can improve efficiency, it might have security implications, such as not adhering strictly to the protocol's design for asynchrony and forward secrecy.

#### **Conclusion**

- **Deviation from Protocol**: Your implementation deviates from the standard Signal Protocol behavior by sending a **Signal message** from **User2** to **User1** without first sending a **PreKey message**.
- **Is It Wrong?**: Technically, it does not adhere to the standard protocol specification. However, whether this is acceptable depends on your application's requirements and security considerations.

### **Possible Reasons in Your Implementation**

#### **1. Session Initialization Logic**

- **User2's Client Logic**: Upon receiving and processing **User1**'s **PreKey message**, **User2**'s client might be setting up a session for both sending and receiving messages.
- **Session Considered Acknowledged**: The client might consider the session fully acknowledged and proceed to send **Signal messages**.

#### **2. PreKey Bundle Availability**

- **Lack of PreKey Bundle Retrieval**: If **User2** does not fetch **User1**'s PreKey bundle before sending the message, it cannot create a **PreKey message**.
- **Assuming Existing Session**: The client might assume an existing session or fall back to sending a **Signal message**.

### **Recommendations**

#### **1. Review Session Establishment Logic**

- **Ensure Directional Sessions**: Verify that your client correctly handles sessions as directional and establishes them independently for sending messages.
- **Modify Client Behavior**: Adjust **User2**'s client to send a **PreKey message** as the first message to **User1**.

#### **2. Fetch PreKey Bundles as Needed**

- **Retrieve User1's PreKey Bundle**: Before **User2** sends the first message to **User1**, it should fetch **User1**'s PreKey bundle to construct a **PreKey message**.
- **Handle PreKey Bundle Storage**: Ensure that the PreKey bundles are stored and managed correctly in your application.

#### **3. Adjust the Test Accordingly**

- **Attempt Decryption as PreKey Message**: After making the necessary adjustments, **User1** should be able to decrypt **User2**'s first message as a **PreKey message**.
- **Verify Session States**: Add debug logs or checks to confirm the session states on both clients during the test.

### **Understanding the Implications**

#### **Security Considerations**

- **Forward Secrecy and Deniability**: Skipping the **PreKey message** might impact these properties.
- **Asynchronous Messaging**: The Signal Protocol is designed to support asynchronous communication securely, and deviating from the protocol might affect this.

#### **Protocol Compliance**

- **Adherence to Specification**: For interoperability and security, it's important to adhere to the protocol specifications.
- **Potential Issues**: Deviations might lead to incompatibilities with other clients or unexpected behaviors.

### **Next Steps**

1. **Investigate User2's Message Sending Logic**

   - Check how **User2** decides whether to send a **PreKey message** or a **Signal message**.
   - Ensure that **User2** fetches **User1**'s PreKey bundle and constructs a **PreKey message** for the first message.

2. **Modify User2's Client Implementation**

   - Implement logic to send a **PreKey message** when no session exists for sending messages to **User1**.
   - Update the session management code to handle sessions correctly in both directions.

3. **Update the Decryption in Your Test**

   - After making the changes, attempt to decrypt **User2**'s first message as a **PreKey message**.
   - Confirm that the decryption succeeds and the message matches the original.

4. **Validate the Entire Message Flow**

   - Run the test to ensure that the message flow now adheres to the expected behavior:
     - **User1 sends PreKey message to User2.**
     - **User2 processes it and sends PreKey message to User1.**
     - **Both users can then exchange Signal messages.**

### **Example Modification to User2's Client**

Ensure that when **User2** sends a message to **User1**, it checks if a session exists:

```go
func (client *ChatClient) SendMessage(ctx context.Context, recipientID string, messageType int, message []byte) error {
    // Check if a session exists for the recipient
    sessionExists := client.SessionStore.HasSession(recipientID)
    if !sessionExists {
        // Fetch recipient's PreKey bundle
        preKeyBundle, err := client.FetchPreKeyBundle(recipientID)
        if err != nil {
            return err
        }
        // Process the PreKey bundle to establish a session
        err = client.Session.ProcessPreKeyBundle(ctx, preKeyBundle)
        if err != nil {
            return err
        }
    }

    // Encrypt and send the message
    encryptedMessage, err := client.Session.EncryptMessage(ctx, message)
    if err != nil {
        return err
    }

    // Send the encrypted message to the recipient
    return client.Transport.SendMessage(recipientID, encryptedMessage)
}
```

### **Adjusting the Test Code**

After making the necessary changes to **User2**'s client, you can attempt to decrypt the message as a **PreKey message**:

```go
// User1 receives and decrypts the PreKey message from User2
select {
case msg2 := <-client1.ChatClient.MessageChannel:
    t.Logf("User1 received message from User2: %v", msg2)
    decryptedMsg2, err := client1.ChatClient.DecryptMessage(context.Background(), msg2, "prekey")
    if err != nil {
        t.Fatalf("User1 failed to decrypt message: %v", err)
    }
    t.Logf("User1 decrypted message: %s", decryptedMsg2)
    if decryptedMsg2 != messageFromUser2 {
        t.Fatalf("User1 decrypted message does not match original message: %s", decryptedMsg2)
    }
```

### **Final Thoughts**

- **Importance of Protocol Compliance**: Adhering to the Signal Protocol's specifications ensures that your application maintains the intended security properties and interoperability.
- **Thorough Testing**: Validate all possible message flows and session states through comprehensive testing to ensure correctness.
- **Security Review**: Consider conducting a security review or audit to identify any deviations or vulnerabilities in your implementation.

### **Conclusion**

- **Explanation**: The reason why decrypting **User2**'s first message as a **Signal message** succeeds is that your implementation is sending a **Signal message** instead of a **PreKey message**.
- **Is It Wrong?**: According to the Signal Protocol, this is not the standard behavior. The first message from **User2** to **User1** should be a **PreKey message** to establish the session in that direction.
- **Action Items**: Review and adjust your client implementation to ensure it sends a **PreKey message** when initiating communication with a user for whom no sending session exists.

---

**Note**: It's crucial for the security and reliability of your chat application to adhere to the Signal Protocol's design principles, especially regarding session establishment and message types. Making these adjustments will help ensure that your application behaves as expected and maintains the highest security standards.