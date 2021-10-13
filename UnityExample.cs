using System;
using System.Collections.Generic;
using proto;
using protoSpace;
using UnityEngine;
using System.Runtime.InteropServices;
using System.Threading;

public class GoDll {
    public static byte[] DefaultBytes = new byte[1024 * 1024];//共享内存
#if UNITY_IPHONE
    const string TECHDLL = "__Internal";
#else
    const string TECHDLL = "godll";
#endif
    [DllImport(TECHDLL, CallingConvention = CallingConvention.Cdecl)]
    public static extern void SetGlobalBytes(byte[] keys, Int32 l);
    [DllImport(TECHDLL, CallingConvention = CallingConvention.Cdecl)]
    public static extern void CallGo(byte[] keys, Int32 l);
    [DllImport(TECHDLL, CallingConvention = CallingConvention.Cdecl)]
    public static extern Int32 GetGoBack();

    public static StartGoDll(){
        SetGlobalBytes(DefaultBytes, DefaultBytes.Length);
        Thread th = new Thread(new ThreadStart(_GetGoBack));
        th.IsBackground = true;
        th.Priority = System.Threading.ThreadPriority.Highest;
        th.Start();
    }

    static _GetGoBack(){
        for(;;){
            try{//阻塞获取buf
                var num = GetGoBack();
                if(num > 0){
                    var caller = System.Text.Encoding.UTF8.GetString(DefaultBytes, num);
                    HandleGoBack(caller);
                }
            }catch(Exception e){
                //err handler
                Debug.LogError(e);
                Thread.Sleep(500);
            }
        }
    }

    public static void HandleGoBack(string caller){
        Debug.Log(caller);
    }
}