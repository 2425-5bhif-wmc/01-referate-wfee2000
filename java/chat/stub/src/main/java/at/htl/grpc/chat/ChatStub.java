package at.htl.grpc.chat;

import io.grpc.Grpc;
import io.grpc.InsecureChannelCredentials;
import io.grpc.ManagedChannel;

public class ChatStub {
    public static void main(String[] args) {
        ManagedChannel channel = Grpc
                .newChannelBuilder("localhost:5555", InsecureChannelCredentials.create())
                .build();
        ChatServiceGrpc.ChatServiceBlockingStub messageStub = ChatServiceGrpc.newBlockingStub(channel);
        String token = messageStub.claimName(Chat.ClaimNameRequest.newBuilder().setName("Winnie").build()).getToken();
    }
}
