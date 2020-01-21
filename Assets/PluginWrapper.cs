using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;
using System.Threading;
using System;

public class PluginWrapper : MonoBehaviour {

  public TextMesh debuglog;
  public GameObject joint;
  public Transform hand;
  public List<Transform> jointList = new List<Transform>();

  public InputField url;

  bool threadRunning;
  Thread pionThread;

  bool pressed = true;
  AndroidJavaClass plugin;

  char[] separator = {' '};
  Int32 count = 45;

  void Start() {
    plugin = new AndroidJavaClass("TestRTC.TestRTC");
    pionThread = new Thread(new ThreadStart(pionThreadFunc));

    // generate hand
    for (int i=0; i<5*3; i++) {
      jointList.Add(Instantiate(joint, hand).transform);
    }
    joint.SetActive(false);
    jointList[3].GetComponent<Renderer>().material.color = new Color(255, 0, 0);
    jointList[2].GetComponent<Renderer>().material.color = new Color(0, 255, 0);
    jointList[1].GetComponent<Renderer>().material.color = new Color(0, 0, 255);
    jointList[0].GetComponent<Renderer>().material.color = new Color(0, 0, 0);
  }

  float divider = 500f;

  void Update() {
    if (OVRInput.Get(OVRInput.Button.PrimaryIndexTrigger)) {
      print("registered");
      print(pressed);
    }
    if (OVRInput.Get(OVRInput.Button.PrimaryIndexTrigger) && pressed) {
      pressed = false;
      print("PRESSED START");
      StartPionThread();
      print("UNPRESSED START");
    }
    if (!pressed) {
      string vars = plugin.CallStatic<string>("returnValues");
      if (vars != "null") {
        String[] strlist = vars.Split(separator, count, StringSplitOptions.RemoveEmptyEntries);
        debuglog.text = strlist.Length.ToString();
        for (int i=0; i<5*3*3; i++) {
          jointList[i].localPosition = new Vector3(-float.Parse(strlist[i*3])/divider, float.Parse(strlist[i*3+1])/divider, float.Parse(strlist[i*3+2])/divider);
        }
      }
    }
  }

  public void StartPionThread() {
    if (!threadRunning) {
      print("STARTING THREAD");
      threadRunning = true;
      pionThread.Start();
      print("THREAD STARTED");
    }
  }

  public void pionThreadFunc() {
    print("ATTATCHING THREAD");
    AndroidJNI.AttachCurrentThread();
    print("THREAD ATTATCHED");
    plugin.CallStatic("startRTC", url.text);
    print("CALLED STARTRTC");
  }

  void OnDisable() {
    if (threadRunning) {
      threadRunning = false;
      pionThread.Join();
      print("THREAD STOPPED");
    }
  }
}
